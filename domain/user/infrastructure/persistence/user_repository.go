package persistence

import (
	"context"
	"food-api/domain/user/domain/model"
	repoDomain "food-api/domain/user/domain/respository"
	"food-api/infrastructure/database"
)

type sqlUserRepo struct {
	Conn *database.Data
}

func NewUserRepository(Conn *database.Data) repoDomain.UserRepository {
	return sqlUserRepo{
		Conn: Conn,
	}
}

func (s sqlUserRepo) GetAllUser(ctx context.Context) ([]model.User, error) {
	panic("implement me")
}

func (s sqlUserRepo) GetById(ctx context.Context, id string) (model.User, error) {
	panic("implement me")
}

func (s sqlUserRepo) Create(ctx context.Context, satellites model.User) error {
	panic("implement me")
}

func (s sqlUserRepo) Update(ctx context.Context, id uint, satellite model.User) error {
	panic("implement me")
}

func (s sqlUserRepo) GetUserByEmailAndPassword(user *model.User) (*model.User, error) {
	panic("implement me")
}
