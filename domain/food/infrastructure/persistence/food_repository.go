package persistence

import (
	"context"
	"food-api/domain/food/application/v1/response"
	"food-api/domain/food/domain/model"
	repoDomain "food-api/domain/food/domain/repository"
	"food-api/infrastructure/database"
	"github.com/google/uuid"
)

type sqlFoodRepo struct {
	Conn *database.Data
}

func NewFoodRepository(Conn *database.Data) repoDomain.FoodRepository {
	return &sqlFoodRepo{
		Conn: Conn,
	}
}

// GetFoodById
func (sr *sqlFoodRepo) GetFoodById(ctx context.Context, id string) (*response.FoodResponse, error) {
	row := sr.Conn.DB.QueryRowContext(ctx, selectFoodById, id)

	var foodResponse response.FoodResponse
	err := row.Scan(&foodResponse.ID, &foodResponse.UserID, &foodResponse.Title, &foodResponse.Description, &foodResponse.FoodImage)
	if err != nil {
		return &response.FoodResponse{}, err
	}

	return &foodResponse, nil
}

// GetFoodByUserId
func (sr *sqlFoodRepo) GetFoodByUserId(ctx context.Context, userId string) (*response.FoodResponse, error) {
	row := sr.Conn.DB.QueryRowContext(ctx, selectFoodByUserId, userId)

	var foodResponse response.FoodResponse
	err := row.Scan(&foodResponse.ID, &foodResponse.UserID, &foodResponse.Title, &foodResponse.Description, &foodResponse.FoodImage)
	if err != nil {
		return &response.FoodResponse{}, err
	}

	return &foodResponse, nil
}

// GetAllFood
func (sr *sqlFoodRepo) GetAllFood(ctx context.Context) ([]response.FoodResponse, error) {
	rows, err := sr.Conn.DB.QueryContext(ctx, selectAllFood)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var foodResponses []response.FoodResponse
	for rows.Next() {
		var foodRow response.FoodResponse
		_ = rows.Scan(&foodRow.ID, &foodRow.UserID, &foodRow.Title, &foodRow.Description, &foodRow.FoodImage)
		foodResponses = append(foodResponses, foodRow)
	}

	return foodResponses, nil
}

// SaveFood
func (sr *sqlFoodRepo) SaveFood(ctx context.Context, food *model.Food) (*response.FoodResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, insertFood)
	if err != nil {
		return &response.FoodResponse{}, err
	}

	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, uuid.New().String(), &food.UserID, &food.Title, &food.Description, &food.FoodImage, &food.CreatedAt, &food.UpdatedAt)

	foodResult := response.FoodResponse{}
	err = row.Scan(&foodResult.ID, &foodResult.UserID, &foodResult.Title, &foodResult.Description, &foodResult.FoodImage)
	if err != nil {
		return &response.FoodResponse{}, err
	}

	return &foodResult, nil
}

// UpdateFood
func (sr *sqlFoodRepo) UpdateFood(ctx context.Context, id string, food *model.Food) error {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, updateFood)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, food.Title, food.Description, food.FoodImage, food.UpdatedAt, id)
	if err != nil {
		return err
	}

	return err
}

// DeleteFood
func (sr *sqlFoodRepo) DeleteFood(ctx context.Context, id string) error {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, deleteFood)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
