package handler

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/widyan/go-http-over-grpc/modules/domain-grpc-test/interfaces"
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

func CreateHandler(usecase interfaces.Usecase_Interface, logger *logrus.Logger, res responses.GinResponses, vldtr validator.ValidatorInterface) *APIHandler {
	return &APIHandler{usecase, logger, res, vldtr}
}

func (a *APIHandler) TestServiceWithParam(ctx context.Context, request *pb.TestRequestLagi) (*pb.TestResponseLagi, error) {
	dcdin, _ := json.Marshal(request)
	a.Logger.Println(dcdin)

	return &pb.TestResponseLagi{
		Status: strconv.FormatInt(int64(request.UserID), 10),
	}, nil
}
