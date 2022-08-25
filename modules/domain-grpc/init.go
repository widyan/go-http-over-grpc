package domain

import (
	"database/sql"
	"log"
	"net"

	"github.com/widyan/go-codebase/responses"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/widyan/go-codebase/modules/domain-grpc/handler"
	pb "github.com/widyan/go-codebase/modules/domain-grpc/proto/v1"
	"github.com/widyan/go-codebase/modules/domain-grpc/repository"
	"github.com/widyan/go-codebase/modules/domain-grpc/usecase"
	validate "github.com/widyan/go-codebase/validator"

	"github.com/sirupsen/logrus"
)

func Init(logger *logrus.Logger, validator validate.ValidatorInterface, pq *sql.DB, cfgRseponses responses.GinResponses) (err error) {

	repo := repository.CreateRepository(pq, pq, logger) // Create transaction from db
	userUsecase := usecase.CreateUsecase(repo, logger)

	// register grpc service server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery())),
		grpc.StreamInterceptor(apmgrpc.NewStreamServerInterceptor()),
	)

	hndlr := handler.CreateHandler(userUsecase, logger, cfgRseponses, validator)
	pb.RegisterTestServer(grpcServer, hndlr)

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return
	}

	reflection.Register(grpcServer)

	log.Println("[SERVER] GRPC server is ready")
	grpcServer.Serve(lis)

	return
}
