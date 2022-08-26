package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	"github.com/widyan/go-codebase/config"
	"github.com/widyan/go-codebase/helper"
	domaingrpc "github.com/widyan/go-codebase/modules/domain-grpc"
	"github.com/widyan/go-codebase/notification"
	pb "github.com/widyan/go-codebase/proto/v1"
	"github.com/widyan/go-codebase/responses"
	"go.elastic.co/apm/module/apmgrpc"
	"go.elastic.co/apm/module/apmhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	validate "github.com/widyan/go-codebase/validator"
)

func main() {
	// Aktivasi environment
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Sprintf("%s: %s", "Failed to load env", err))
	}

	response := responses.CreateCustomResponses(os.Getenv("PROJECT_NAME"))

	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			return funcname, f.File + ":" + strconv.Itoa(f.Line)
		},
	})

	logger.SetReportCaller(true)
	// logger.Hooks.Add(&apmlogrus.Hook{
	// 	LogLevels: logrus.AllLevels,
	// })

	// Get dependency call API
	toolsAPI := helper.CreateToolsAPI(logger)

	// ************************************ Setting notification to telegram if become error ***********
	notifTelegram := notification.CreateNotification(toolsAPI, os.Getenv("PROJECT_NAME"), os.Getenv("TOKEN_BOT_TELEGRAM"), os.Getenv("CHAT_ID"))
	logger.AddHook(notifTelegram)
	// *************************************************************************************************

	// ************************************ Get dependency validator ***********************************
	validator := validator.New()
	vldt := validate.CreateValidator(validator)
	// *************************************************************************************************

	var cfg config.Config

	/*
		// ************************************ Config for implement DB ************************************
		pq := cfg.Postgresql(os.Getenv("GORM_CONNECTION"), 20, 20)
		pqdbAuth := cfg.Postgresql(os.Getenv("POSTGRES_AUTH_CONNECTION"), 20, 20)
		// *************************************************************************************************

		auth := middleware.Init(routesGin, logger, pqdbAuth, vldt, response)
		domain.Init(routesGin, logger, vldt, pq, response, auth)
	*/

	// create run group
	g, _ := errgroup.WithContext(context.Background())

	var servers []*http.Server
	// goroutine to check for signals to gracefully finish all functions
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			log.Println("received signal: %s\n", sig)
			logger.Println("Close connection postgresql")

			for i, s := range servers {
				if err := s.Shutdown(context.Background()); err != nil {
					if err == nil {
						log.Println("error shutting down server %d: %v", i, err)
						panic(err)
					}
				}
			}
			os.Exit(1)
		}
		return nil
	})

	g.Go(func() error { return NewGRPCServer(logger, cfg, vldt, response) })
	g.Go(func() error { return NewHTTPServer(logger) })

	if err := g.Wait(); !IgnoreErr(err) {
		panic(err)
	}

}

func NewGRPCServer(logger *logrus.Logger, cfg config.Config, validator validate.ValidatorInterface, cfgRseponses responses.GinResponses) (err error) {
	// register grpc service server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery())),
		grpc.StreamInterceptor(apmgrpc.NewStreamServerInterceptor()),
	)

	if err = domaingrpc.Init(logger, grpcServer, cfg, validator, cfgRseponses); err != nil {
		return
	}

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return
	}

	reflection.Register(grpcServer)

	logger.Println("[SERVER] GRPC server is ready")
	grpcServer.Serve(lis)
	return nil
}

func NewHTTPServer(logger *logrus.Logger) error {
	/*
		tracer.Start(tracer.WithDebugMode(true))
		defer tracer.Stop()
	*/

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	addr := "0.0.0.0:5000"
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()),
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
	)
	if err != nil {
		logger.Fatalf("fail to dial: %v", err)
		panic(err)
	}
	defer conn.Close()

	// Create new grpc-gateway
	rmux := gwruntime.NewServeMux(gwruntime.WithMarshalerOption(gwruntime.MIMEWildcard, &gwruntime.JSONPb{
		EmitDefaults: true,
		OrigName:     true,
	}))

	// register gateway endpoints
	for _, f := range []func(ctx context.Context, mux *gwruntime.ServeMux, conn *grpc.ClientConn) error{
		// register grpc service handler
		pb.RegisterTestHandler,
	} {
		if err = f(ctx, rmux, conn); err != nil {
			log.Fatal(err)
			panic(err)

		}
	}

	// logrus configuration
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	// create http server mux
	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.Handle("/metrics", promhttp.Handler())
	apmhttp.Wrap(mux)

	/*
		// run swagger server
		if hlp.Getenv("DEVELOPMENT") == "true" {
			NewHTTPEncodedAndSwaggerServer(mux)
		}
	*/

	// running delivery http server
	logger.Println("[SERVER] REST HTTP server is ready")

	if err = http.ListenAndServe(":7000", mux); err != nil {
		logger.Println(err.Error())
		return err
	}

	return nil
}

// ignoreErr returns true when err can be safely ignored.
func IgnoreErr(err error) bool {
	switch {
	case err == nil || err == http.ErrServerClosed || err == cmux.ErrListenerClosed:
		return true
	}
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err.Error() == "use of closed network connection"
	}
	return false
}
