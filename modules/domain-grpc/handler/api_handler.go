package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/widyan/go-codebase/modules/domain-grpc/interfaces"
	pb "github.com/widyan/go-codebase/proto/v1"
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
	dcdin, _ := json.Marshal(request)
	fmt.Println(string(dcdin))

	return &pb.TestResponse{
		Status: request.Name,
	}, nil
}

func (a *APIHandler) TestServiceWithParam(ctx context.Context, request *pb.TestRequest) (*pb.TestResponse, error) {
	dcdin, _ := json.Marshal(request)
	fmt.Println(string(dcdin))

	return &pb.TestResponse{
		Status: strconv.FormatInt(int64(request.UserID), 10),
	}, nil
}
