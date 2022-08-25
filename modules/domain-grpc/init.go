package domain

import (
	"github.com/widyan/go-codebase/config"
	"github.com/widyan/go-codebase/responses"
	"google.golang.org/grpc"

	"github.com/widyan/go-codebase/modules/domain-grpc/handler"
	pb "github.com/widyan/go-codebase/modules/domain-grpc/proto/v1"
	"github.com/widyan/go-codebase/modules/domain-grpc/repository"
	"github.com/widyan/go-codebase/modules/domain-grpc/usecase"
	validate "github.com/widyan/go-codebase/validator"

	"github.com/sirupsen/logrus"
)

func Init(logger *logrus.Logger, grpcServer *grpc.Server, cfg config.Config, validator validate.ValidatorInterface, cfgRseponses responses.GinResponses) (err error) {
	repo := repository.CreateRepository(nil, nil, logger) // Create transaction from db
	userUsecase := usecase.CreateUsecase(repo, logger)

	hndlr := handler.CreateHandler(userUsecase, logger, cfgRseponses, validator)
	pb.RegisterTestServer(grpcServer, hndlr)

	return
}
