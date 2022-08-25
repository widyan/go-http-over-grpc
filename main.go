package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/soheilhy/cmux"
	"github.com/widyan/go-codebase/helper"
	initial "github.com/widyan/go-codebase/modules/domain-grpc"
	pb "github.com/widyan/go-codebase/modules/domain-grpc/proto/v1"
	"github.com/widyan/go-codebase/notification"
	"github.com/widyan/go-codebase/responses"
	"go.elastic.co/apm/module/apmgrpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

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
			_, filename := path.Split(f.File)
			return funcname, filename
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

	// var cfg config.Config

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

	g.Go(func() error { return initial.Init(logger, vldt, nil, response) })
	g.Go(func() error { return NewHTTPServer() })

	if err := g.Wait(); !IgnoreErr(err) {
		panic(err)
	}

}

func NewGRPCServer() error {
	return nil
}

func NewHTTPServer() error {
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
		log.Fatalf("fail to dial: %v", err)
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
	// mux.Handle("/metrics", promhttp.Handler())
	// apmhttp.Wrap(mux)

	/*
		// run swagger server
		if hlp.Getenv("DEVELOPMENT") == "true" {
			NewHTTPEncodedAndSwaggerServer(mux)
		}
	*/

	// running delivery http server
	log.Println("[SERVER] REST HTTP server is ready")

	http.ListenAndServe(":7000", mux)

	if err != nil {
		log.Println(err.Error())
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
