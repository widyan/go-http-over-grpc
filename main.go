package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/widyan/go-http-over-grpc/config"
	"github.com/widyan/go-http-over-grpc/gateway"
	"github.com/widyan/go-http-over-grpc/helper"
	"github.com/widyan/go-http-over-grpc/insecure"
	domaingrpc "github.com/widyan/go-http-over-grpc/modules/domain-grpc"
	domaingrpclagi "github.com/widyan/go-http-over-grpc/modules/domain-grpc-test"
	"github.com/widyan/go-http-over-grpc/notification"
	"github.com/widyan/go-http-over-grpc/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	validate "github.com/widyan/go-http-over-grpc/validator"
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

	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	addr := "0.0.0.0:10000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Panic("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer(
		// TODO: Replace with your own certificate!
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
	)

	domaingrpc.Init(logger, grpcServer, cfg, vldt, response)
	domaingrpclagi.Init(logger, grpcServer, cfg, vldt, response)

	// Serve gRPC Server
	log.Info("Serving gRPC on https://", addr)
	go func() {
		logger.Panic(grpcServer.Serve(lis))
	}()

	err = gateway.Run(logger, "dns:///"+addr)
	logger.Panic(err) //
}
