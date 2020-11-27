package respository

import (
	"context"
	"food-api/domain/food/domain/model"
)

type FoodRepository interface {
	SaveFood(ctx context.Context, food *model.Food) (*model.Food, error)
	GetFood(ctx context.Context, id string) (*model.Food, error)
	GetAllFood(ctx context.Context) ([]model.Food, error)
	UpdateFood(ctx context.Context, food *model.Food) (*model.Food, error)
	DeleteFood(ctx context.Context, id string) error
}
