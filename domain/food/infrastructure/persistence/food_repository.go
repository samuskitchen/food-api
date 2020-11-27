package persistence

import (
	"context"
	"food-api/domain/food/application/v1/response"
	"food-api/domain/food/domain/model"
	repoDomain "food-api/domain/food/domain/respository"
	"food-api/infrastructure/database"
	"time"
)

type sqlFoodRepo struct {
	Conn *database.Data
}

func NewFoodRepository(Conn *database.Data) repoDomain.FoodRepository {
	return &sqlFoodRepo{
		Conn: Conn,
	}
}

func (sr *sqlFoodRepo) SaveFood(ctx context.Context, food *model.Food) (*response.FoodResponse, error) {
	now := time.Now() //.Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	stmt, err := sr.Conn.DB.PrepareContext(ctx, insertFood)
	if err != nil {
		return &response.FoodResponse{}, err
	}

	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, food.UserID, food.Title, food.Description, food.FoodImage, now, now)

	foodResponse := response.FoodResponse{}
	err = row.Scan(&foodResponse.ID, &foodResponse.UserID, &foodResponse.Title, &foodResponse.Description, &foodResponse.FoodImage)
	if err != nil {
		return &response.FoodResponse{}, err
	}

	return &foodResponse, nil
}

func (sr *sqlFoodRepo) GetFoodById(ctx context.Context, id string) (*response.FoodResponse, error) {
	row := sr.Conn.DB.QueryRowContext(ctx, selectFoodById, id)

	var foodResponse response.FoodResponse
	err := row.Scan(&foodResponse.ID, &foodResponse.UserID, &foodResponse.Title, &foodResponse.Description, &foodResponse.FoodImage)
	if err != nil {
		return &response.FoodResponse{}, err
	}

	return &foodResponse, nil
}

func (sr *sqlFoodRepo) GetFoodByUserId(ctx context.Context, userId string) (*response.FoodResponse, error) {
	row := sr.Conn.DB.QueryRowContext(ctx, selectFoodByUserId, userId)

	var foodResponse response.FoodResponse
	err := row.Scan(&foodResponse.ID, &foodResponse.UserID, &foodResponse.Title, &foodResponse.Description, &foodResponse.FoodImage)
	if err != nil {
		return &response.FoodResponse{}, err
	}

	return &foodResponse, nil
}

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

func (sr *sqlFoodRepo) UpdateFood(ctx context.Context, id string, food *model.Food) error {
	now := time.Now() //.Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	stmt, err := sr.Conn.DB.PrepareContext(ctx, updateFood)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, food.Title, food.Description, food.FoodImage, now, id)
	if err != nil {
		return err
	}

	return err
}

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