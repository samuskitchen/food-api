package repository

import (
	"context"
	"food-api/domain/food/application/v1/response"
	"food-api/domain/food/domain/model"
)

type FoodRepository interface {
	SaveFood(ctx context.Context, food *model.Food) (*response.FoodResponse, error)
	GetFoodById(ctx context.Context, id string) (*response.FoodResponse, error)
	GetFoodByUserId(ctx context.Context, id string) (*response.FoodResponse, error)
	GetAllFood(ctx context.Context) ([]response.FoodResponse, error)
	UpdateFood(ctx context.Context, id string, food *model.Food) error
	DeleteFood(ctx context.Context, id string) error
}