package v1

import (
	repoDomain "food-api/domain/food/domain/respository"
	"food-api/domain/food/infrastructure/persistence"
	"food-api/infrastructure/database"
)

// FoodRouter
type FoodRouter struct {
	Repo repoDomain.FoodRepository
}

func NewFoodHandler(db *database.Data) *FoodRouter {
	return &FoodRouter{
		Repo: persistence.NewFoodRepository(db),
	}
}

