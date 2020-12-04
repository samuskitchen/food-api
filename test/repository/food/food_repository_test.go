package food

import (
	"context"
	"database/sql"
	responseFood "food-api/domain/food/application/v1/response"
	"food-api/domain/food/domain/model"
	"food-api/domain/food/domain/repository"
	"food-api/domain/food/infrastructure/persistence"
	"food-api/infrastructure/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

// represent the repository
var (
	dbMockFoods        *sql.DB
	connMockFood       database.Data
	foodRepositoryMock repository.FoodRepository
)

// NewMockFood initialize mock connection to database
func NewMockFood() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbMockFoods = db
	connMockFood = database.Data{
		DB: dbMockFoods,
	}

	foodRepositoryMock = persistence.NewFoodRepository(&connMockFood)

	return mock
}

// CloseMockFood Close attaches the provider and close the connection
func CloseMockFood() {
	err := dbMockFoods.Close()
	if err != nil {
		log.Println("Error close database test")
	}
}

// dataFood is data for test
func dataFood() []model.Food {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)
	userId := uuid.New().String()

	return []model.Food{
		{
			ID:          uuid.New().String(),
			UserID:      userId,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New().String(),
			UserID:      userId,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}

// dataFoodResponse is data for test
func dataFoodResponse() []responseFood.FoodResponse {
	userId := uuid.New().String()

	return []responseFood.FoodResponse{
		{
			ID:          uuid.New().String(),
			UserID:      userId,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
		},
		{
			ID:          uuid.New().String(),
			UserID:      userId,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
		},
	}
}

func Test_sqlFoodRepo_GetAllFood(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		mock.ExpectQuery("SELECT 1 FROM user")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		users, err := foodRepositoryMock.GetAllFood(ctx)
		assert.Error(tt, err)
		assert.Nil(tt, users)
	})

	t.Run("Get All Food Successful", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		foodsData := dataFoodResponse()
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "food_image"}).
			AddRow(foodsData[0].ID, foodsData[0].UserID, foodsData[0].Title, foodsData[0].Description, foodsData[0].FoodImage).
			AddRow(foodsData[1].ID, foodsData[0].UserID, foodsData[1].Title, foodsData[1].Description, foodsData[1].FoodImage)

		mock.ExpectQuery(selectAllFoodTest).WillReturnRows(rows)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		foods, err := foodRepositoryMock.GetAllFood(ctx)
		assert.NotEmpty(tt, foods)
		assert.NoError(tt, err)
		assert.Len(tt, foods, 2)
	})

}

func Test_sqlFoodRepo_GetFoodById(t *testing.T) {
	foodTest := dataFoodResponse()[0]

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		row := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "food_image"}).
			AddRow(foodTest.ID, foodTest.UserID, foodTest.Title, foodTest.Description, foodTest.FoodImage)

		mock.ExpectQuery(selectFoodByIdTest).WithArgs(nil).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		foodResult, err := foodRepositoryMock.GetFoodById(ctx, foodTest.ID)
		assert.Error(tt, err)
		assert.NotNil(tt, foodResult)
	})

	t.Run("Get Food By Id Successful", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()
		row := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "food_image"}).
			AddRow(foodTest.ID, foodTest.UserID, foodTest.Title, foodTest.Description, foodTest.FoodImage)

		mock.ExpectQuery(selectFoodByIdTest).WithArgs(foodTest.ID).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		foodResult, err := foodRepositoryMock.GetFoodById(ctx, foodTest.ID)
		assert.NoError(tt, err)
		assert.NotNil(tt, foodResult)
	})

}

func Test_sqlFoodRepo_GetFoodByUserId(t *testing.T) {
	foodTest := dataFoodResponse()[0]

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		row := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "food_image"}).
			AddRow(foodTest.ID, foodTest.UserID, foodTest.Title, foodTest.Description, foodTest.FoodImage)

		mock.ExpectQuery(selectFoodByUserIdTest).WithArgs(nil).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		foodResult, err := foodRepositoryMock.GetFoodByUserId(ctx, foodTest.ID)
		assert.Error(tt, err)
		assert.NotNil(tt, foodResult)
	})

	t.Run("Get Food By User Successful", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()
		row := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "food_image"}).
			AddRow(foodTest.ID, foodTest.UserID, foodTest.Title, foodTest.Description, foodTest.FoodImage)

		mock.ExpectQuery(selectFoodByUserIdTest).WithArgs(foodTest.UserID).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		foodResult, err := foodRepositoryMock.GetFoodByUserId(ctx, foodTest.UserID)
		assert.NoError(tt, err)
		assert.NotNil(tt, foodResult)
	})
}

func Test_sqlFoodRepo_SaveFood(t *testing.T) {

}

func Test_sqlFoodRepo_UpdateFood(t *testing.T) {

}

func Test_sqlFoodRepo_DeleteFood(t *testing.T) {

}
