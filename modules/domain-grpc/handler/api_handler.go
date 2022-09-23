package handler

import (
	"context"
	"fmt"

	"github.com/widyan/go-http-over-grpc/modules/domain-grpc/interfaces"
	pb "github.com/widyan/go-http-over-grpc/proto/latest"
	"github.com/widyan/go-http-over-grpc/responses"
	"github.com/widyan/go-http-over-grpc/validator"

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
	/*
		dcdin, _ := json.Marshal(request)
		fmt.Println(string(dcdin))
	*/

	fmt.Println(request.Validate())

	return &pb.TestResponse{
		Status: request.UserID,
	}, nil
}

func (a *APIHandler) TestServiceWithParam(ctx context.Context, request *pb.TestRequest) (*pb.TestResponse, error) {
	/*
		dcdin, _ := json.Marshal(request)
		fmt.Println(string(dcdin))
	*/

	return &pb.TestResponse{
		Status: request.UserID,
	}, nil
}
