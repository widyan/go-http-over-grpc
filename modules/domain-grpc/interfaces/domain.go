package interfaces

import (
	"context"

	"github.com/widyan/go-codebase/modules/domain-grpc/entity"
)

type Usecase_Interface interface {
	Test(ctx context.Context, name string) string
}

type Repository_Interface interface {
	InsertUser(ctx context.Context, user entity.Users) (err error)
	GetOneUser(ctx context.Context) (user entity.Users, err error)
	GetAllUsers(ctx context.Context) (users []entity.Users, err error)
	UpdateUserByID(ctx context.Context, id int, fullname string) (err error)
	GetOneUserByID(ctx context.Context, id int) (user entity.Users, err error)
}

type Scheduller interface {
	TestScheduller()
}
