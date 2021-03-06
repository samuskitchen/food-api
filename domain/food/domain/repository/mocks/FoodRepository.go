// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"
	model "food-api/domain/food/domain/model"

	mock "github.com/stretchr/testify/mock"

	response "food-api/domain/food/application/v1/response"
)

// FoodRepository is an autogenerated mock type for the FoodRepository type
type FoodRepository struct {
	mock.Mock
}

// DeleteFood provides a mock function with given fields: ctx, id
func (_m *FoodRepository) DeleteFood(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllFood provides a mock function with given fields: ctx
func (_m *FoodRepository) GetAllFood(ctx context.Context) ([]response.FoodResponse, error) {
	ret := _m.Called(ctx)

	var r0 []response.FoodResponse
	if rf, ok := ret.Get(0).(func(context.Context) []response.FoodResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]response.FoodResponse)
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

// GetFoodById provides a mock function with given fields: ctx, id
func (_m *FoodRepository) GetFoodById(ctx context.Context, id string) (*response.FoodResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 *response.FoodResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) *response.FoodResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.FoodResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFoodByUserId provides a mock function with given fields: ctx, id
func (_m *FoodRepository) GetFoodByUserId(ctx context.Context, id string) (*response.FoodResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 *response.FoodResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) *response.FoodResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.FoodResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveFood provides a mock function with given fields: ctx, food
func (_m *FoodRepository) SaveFood(ctx context.Context, food *model.Food) (*response.FoodResponse, error) {
	ret := _m.Called(ctx, food)

	var r0 *response.FoodResponse
	if rf, ok := ret.Get(0).(func(context.Context, *model.Food) *response.FoodResponse); ok {
		r0 = rf(ctx, food)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.FoodResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.Food) error); ok {
		r1 = rf(ctx, food)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFood provides a mock function with given fields: ctx, id, food
func (_m *FoodRepository) UpdateFood(ctx context.Context, id string, food *model.Food) error {
	ret := _m.Called(ctx, id, food)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Food) error); ok {
		r0 = rf(ctx, id, food)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
