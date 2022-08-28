package domain

import (
	"github.com/widyan/go-http-over-grpc/config"
	"github.com/widyan/go-http-over-grpc/responses"
	"google.golang.org/grpc"

	"github.com/widyan/go-http-over-grpc/modules/domain-grpc/handler"
	"github.com/widyan/go-http-over-grpc/modules/domain-grpc/repository"
	"github.com/widyan/go-http-over-grpc/modules/domain-grpc/usecase"
	pb "github.com/widyan/go-http-over-grpc/proto/latest"
	validate "github.com/widyan/go-http-over-grpc/validator"

	"github.com/sirupsen/logrus"
)

func Init(logger *logrus.Logger, grpcServer *grpc.Server, cfg config.Config, validator validate.ValidatorInterface, cfgRseponses responses.GinResponses) {
	repo := repository.CreateRepository(nil, nil, logger) // Create transaction from db
	userUsecase := usecase.CreateUsecase(repo, logger)

	hndlr := handler.CreateHandler(userUsecase, logger, cfgRseponses, validator)
	pb.RegisterTestServer(grpcServer, hndlr)
}
