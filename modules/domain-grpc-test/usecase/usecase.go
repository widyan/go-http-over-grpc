package usecase

import (
	"context"

	"github.com/widyan/go-http-over-grpc/modules/domain-grpc-test/interfaces"

	"github.com/sirupsen/logrus"
)

type Usecase struct {
	Repository interfaces.Repository_Interface
	Logger     *logrus.Logger
}

func CreateUsecase(repo interfaces.Repository_Interface, logger *logrus.Logger) interfaces.Usecase_Interface {
	return &Usecase{
		Repository: repo,
		Logger:     logger,
	}
}

func (b *Usecase) Test(ctx context.Context, name string) string {
	return name
}
