package persistence

import (
	"context"
	"food-api/domain/food/domain/model"
	repoDomain "food-api/domain/food/domain/respository"
	"food-api/infrastructure/database"
)

type sqlFoodRepo struct {
	Conn *database.Data
}

func NewFoodRepository(Conn *database.Data) repoDomain.FoodRepository{
	return &sqlFoodRepo{
		Conn: Conn,
	}
}

func (s sqlFoodRepo) SaveFood(ctx context.Context, food *model.Food) (*model.Food, error) {
	panic("implement me")
}

func (s sqlFoodRepo) GetFood(ctx context.Context, id string) (*model.Food, error) {
	panic("implement me")
}

func (s sqlFoodRepo) GetAllFood(ctx context.Context) ([]model.Food, error) {
	panic("implement me")
}

func (s sqlFoodRepo) UpdateFood(ctx context.Context, food *model.Food) (*model.Food, error) {
	panic("implement me")
}

func (s sqlFoodRepo) DeleteFood(ctx context.Context, id string) error {
	panic("implement me")
}
