// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/widyan/go-codebase/modules/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// Repository_Interface is an autogenerated mock type for the Repository_Interface type
type Repository_Interface struct {
	mock.Mock
}

// GetAllUsers provides a mock function with given fields: ctx
func (_m *Repository_Interface) GetAllUsers(ctx context.Context) ([]entity.Users, error) {
	ret := _m.Called(ctx)

	var r0 []entity.Users
	if rf, ok := ret.Get(0).(func(context.Context) []entity.Users); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Users)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOneUser provides a mock function with given fields: ctx
func (_m *Repository_Interface) GetOneUser(ctx context.Context) (entity.Users, error) {
	ret := _m.Called(ctx)

	var r0 entity.Users
	if rf, ok := ret.Get(0).(func(context.Context) entity.Users); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(entity.Users)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOneUserByID provides a mock function with given fields: ctx, id
func (_m *Repository_Interface) GetOneUserByID(ctx context.Context, id int) (entity.Users, error) {
	ret := _m.Called(ctx, id)

	var r0 entity.Users
	if rf, ok := ret.Get(0).(func(context.Context, int) entity.Users); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(entity.Users)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertUser provides a mock function with given fields: ctx, user
func (_m *Repository_Interface) InsertUser(ctx context.Context, user entity.Users) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Users) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUserByID provides a mock function with given fields: ctx, id, fullname
func (_m *Repository_Interface) UpdateUserByID(ctx context.Context, id int, fullname string) error {
	ret := _m.Called(ctx, id, fullname)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, id, fullname)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewRepository_InterfaceT interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository_Interface creates a new instance of Repository_Interface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository_Interface(t NewRepository_InterfaceT) *Repository_Interface {
	mock := &Repository_Interface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
