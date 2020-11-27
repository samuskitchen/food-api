package respository

import (
	"context"
	"food-api/domain/user/application/v1/response"
	"food-api/domain/user/domain/model"
)

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]response.UserResponse, error)
	GetById(ctx context.Context, id string) (response.UserResponse, error)
	CreateUser(ctx context.Context, user model.User) (response.UserResponse, error)
	UpdateUser(ctx context.Context, id string, user model.User) error
	GetUserByEmailAndPassword(ctx context.Context, user *model.User) (*response.UserResponse, error)
}