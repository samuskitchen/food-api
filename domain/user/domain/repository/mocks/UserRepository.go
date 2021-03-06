// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"
	model "food-api/domain/user/domain/model"

	mock "github.com/stretchr/testify/mock"

	response "food-api/domain/user/application/v1/response"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *UserRepository) CreateUser(ctx context.Context, user *model.User) (*response.UserResponse, error) {
	ret := _m.Called(ctx, user)

	var r0 *response.UserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) *response.UserResponse); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.UserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllUser provides a mock function with given fields: ctx
func (_m *UserRepository) GetAllUser(ctx context.Context) ([]response.UserResponse, error) {
	ret := _m.Called(ctx)

	var r0 []response.UserResponse
	if rf, ok := ret.Get(0).(func(context.Context) []response.UserResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]response.UserResponse)
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

// GetById provides a mock function with given fields: ctx, id
func (_m *UserRepository) GetById(ctx context.Context, id string) (response.UserResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 response.UserResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) response.UserResponse); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(response.UserResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmailAndPassword provides a mock function with given fields: ctx, user
func (_m *UserRepository) GetUserByEmailAndPassword(ctx context.Context, user *model.User) (*response.UserResponse, error) {
	ret := _m.Called(ctx, user)

	var r0 *response.UserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) *response.UserResponse); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.UserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, id, user
func (_m *UserRepository) UpdateUser(ctx context.Context, id string, user model.User) error {
	ret := _m.Called(ctx, id, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.User) error); ok {
		r0 = rf(ctx, id, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
