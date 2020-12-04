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

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		prep := mock.ExpectPrepare("insertFoodTest")
		prep.ExpectExec().
			WithArgs(dataFood()[0].ID, dataFood()[0].UserID, dataFood()[0].Title, dataFood()[0].Description, dataFood()[0].FoodImage, dataFood()[0].CreatedAt, dataFood()[0].UpdatedAt).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		foodResult, err := foodRepositoryMock.SaveFood(ctx, &dataFood()[0])
		assert.Error(tt, err)
		assert.NotNil(tt, foodResult)
	})

	t.Run("Error Scan Row", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		prep := mock.ExpectPrepare(insertFoodTest)
		prep.ExpectQuery().
			WithArgs(dataFood()[0].ID, dataFood()[0].UserID, dataFood()[0].Title, dataFood()[0].Description, dataFood()[0].FoodImage, dataFood()[0].CreatedAt, dataFood()[0].UpdatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"first_name"}).AddRow("Error"))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		foodResult, err := foodRepositoryMock.SaveFood(ctx, &dataFood()[0])
		assert.Error(tt, err)
		assert.NotNil(tt, foodResult)
	})

	t.Run("Create Food Successful", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		dataTest := dataFood()[0]
		prep := mock.ExpectPrepare(insertFoodTest)
		prep.ExpectQuery().
			WithArgs(dataTest.ID, dataTest.UserID, dataTest.Title, dataTest.Description, dataTest.FoodImage, dataTest.CreatedAt, dataTest.UpdatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "description", "food_image"}).
				AddRow(dataTest.ID, dataTest.UserID, dataTest.Title, dataTest.Description, dataTest.FoodImage))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		foodResult, err := foodRepositoryMock.SaveFood(ctx, &dataTest)
		assert.NoError(tt, err)
		assert.NotNil(tt, foodResult)
	})
}

func Test_sqlFoodRepo_UpdateFood(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		prep := mock.ExpectPrepare("updateFoodTest")
		prep.ExpectExec().
			WithArgs(dataFood()[0].Title, dataFood()[0].Description, dataFood()[0].FoodImage, dataFood()[0].UpdatedAt, dataFood()[0].ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := foodRepositoryMock.UpdateFood(ctx, dataFood()[0].ID, &dataFood()[0])
		assert.Error(tt, err)
	})

	t.Run("Error Statement SQL", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		prep := mock.ExpectPrepare(updateFoodTest)
		prep.ExpectExec().
			WithArgs(dataFood()[0].Title, dataFood()[0].Description, dataFood()[0].FoodImage, dataFood()[0].UpdatedAt, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := foodRepositoryMock.UpdateFood(ctx, dataFood()[0].ID, &dataFood()[0])
		assert.Error(tt, err)
	})

	t.Run("Update Food Successful", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		dataTest := dataFood()[0]
		prep := mock.ExpectPrepare(updateFoodTest)
		prep.ExpectExec().
			WithArgs(dataTest.Title, dataTest.Description, dataTest.FoodImage, dataTest.UpdatedAt, dataTest.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := foodRepositoryMock.UpdateFood(ctx, dataTest.ID, &dataTest)
		assert.NoError(tt, err)
	})

}

func Test_sqlFoodRepo_DeleteFood(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		prep := mock.ExpectPrepare("deleteFoodTest")
		prep.ExpectExec().
			WithArgs(uint(1)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := foodRepositoryMock.DeleteFood(ctx, "")
		assert.Error(tt, err)
	})

	t.Run("Error Statement SQL", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		prep := mock.ExpectPrepare(deleteFoodTest)
		prep.ExpectExec().
			WithArgs(nil).
			WillReturnResult(sqlmock.NewResult(0, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := foodRepositoryMock.DeleteFood(ctx, "")
		assert.Error(tt, err)
	})

	t.Run("Delete Food Successful", func(tt *testing.T) {
		mock := NewMockFood()
		defer func() {
			CloseMockFood()
		}()

		userID := dataFood()[0].ID

		prep := mock.ExpectPrepare(deleteFoodTest)
		prep.ExpectExec().
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := foodRepositoryMock.DeleteFood(ctx, userID)
		assert.NoError(tt, err)
	})

}
