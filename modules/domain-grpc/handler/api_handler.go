package handler

import (
	"context"

	"github.com/widyan/go-codebase/modules/domain-grpc/interfaces"
	pb "github.com/widyan/go-codebase/modules/domain-grpc/proto/v1"
	"github.com/widyan/go-codebase/responses"
	"github.com/widyan/go-codebase/validator"

	"github.com/sirupsen/logrus"
)

type APIHandler struct {
	Usecase   interfaces.Usecase_Interface
	Logger    *logrus.Logger
	Res       responses.GinResponses
	Validator validator.ValidatorInterface
}

var usecase interfaces.Usecase_Interface
var customLogger *logrus.Logger
var response responses.GinResponses
var validate validator.ValidatorInterface

func CreateHandler(Usecase interfaces.Usecase_Interface, logger *logrus.Logger, res responses.GinResponses, vldtr validator.ValidatorInterface) *APIHandler {
	return &APIHandler{usecase, customLogger, response, validate}
}

func (a *APIHandler) TestService(ctx context.Context, request *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{
		Status: "Ok",
	}, nil
}
