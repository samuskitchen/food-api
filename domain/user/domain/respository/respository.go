package respository

import (
	"context"
	"food-api/domain/user/domain/model"
)

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]model.User, error)
	GetById(ctx context.Context, id string) (model.User, error)
	Create(ctx context.Context, satellites model.User) error
	Update(ctx context.Context, id uint, satellite model.User) error
	GetUserByEmailAndPassword(*model.User) (*model.User, error)
}
