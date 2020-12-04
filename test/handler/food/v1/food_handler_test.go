package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	v1 "food-api/domain/food/application/v1"
	responseFood "food-api/domain/food/application/v1/response"
	"food-api/domain/food/domain/model"
	repoMock "food-api/domain/food/domain/repository/mocks"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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

func TestFoodRouter_GetAllFood(t *testing.T) {

	t.Run("Error Get All Food Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/api/v1/foods/", nil)
		response := httptest.NewRecorder()
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("GetAllFood", mock.Anything, mock.Anything).Return(nil, errors.New("error trace test"))

		testFoodHandler.GetAllFood(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get All Food Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/api/v1/foods/", nil)
		response := httptest.NewRecorder()
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("GetAllFood", mock.Anything).Return(dataFoodResponse(), nil)

		testFoodHandler.GetAllFood(response, request)
		mockRepository.AssertExpectations(tt)
	})

}

func TestFoodRouter_GetOneByUserHandler(t *testing.T) {

	t.Run("Error Param Get One By User Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/foods/user/{id}", nil)

		mockRepository := &repoMock.FoodRepository{}
		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}

		testFoodHandler.GetOneByUserHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Get One By User Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/foods/user/{id}", nil)

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("GetFoodByUserId", mock.Anything, mock.Anything).Return(&responseFood.FoodResponse{}, errors.New("error sql")).Once()

		testFoodHandler.GetOneByUserHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get One By User Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/foods/user/{id}", nil)

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("GetFoodByUserId", mock.Anything, mock.Anything).Return(&dataFoodResponse()[0], nil).Once()

		testFoodHandler.GetOneByUserHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})
}

func TestFoodRouter_GetOneHandler(t *testing.T) {

	t.Run("Error Param Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/foods/{id}", nil)

		mockRepository := &repoMock.FoodRepository{}
		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}

		testFoodHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/foods/{id}", nil)

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("GetFoodById", mock.Anything, mock.Anything).Return(&responseFood.FoodResponse{}, errors.New("error sql")).Once()

		testFoodHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/foods/{id}", nil)

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("GetFoodById", mock.Anything, mock.Anything).Return(&dataFoodResponse()[0], nil).Once()

		testFoodHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})
}

func TestFoodRouter_CreateHandler(t *testing.T) {

	t.Run("Error Body Create Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPost, "/api/v1/foods/", bytes.NewReader(nil))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}

		testFoodHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Create Handler", func(tt *testing.T) {
		dataFood()[0].ID = ""

		marshal, err := json.Marshal(dataFood()[0])
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/foods/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("SaveFood", mock.Anything, mock.Anything).Return(&responseFood.FoodResponse{}, errors.New("error sql"))

		testFoodHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Validate Create Handler", func(tt *testing.T) {

		var userTest = dataFood()[0]
		userTest.Title = ""

		marshal, err := json.Marshal(userTest)
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/foods/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}

		testFoodHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Create Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataFood()[0])
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/foods/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("SaveFood", mock.Anything, mock.Anything).Return(&dataFoodResponse()[0], nil)

		testFoodHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

}

func TestFoodRouter_UpdateHandler(t *testing.T) {

	t.Run("Error Param Update Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/api/v1/foods/{id}", nil)

		mockRepository := &repoMock.FoodRepository{}
		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}

		testFoodHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Body Update Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPut, "/api/v1/foods/{id}", bytes.NewReader(nil))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))

		mockRepository := &repoMock.FoodRepository{}
		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}

		testFoodHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Update Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataFood()[0])
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/foods/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("UpdateFood", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error sql")).Once()

		testFoodHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Validate Update Handler", func(tt *testing.T) {

		var userTest = dataFood()[0]
		userTest.Title = ""

		marshal, err := json.Marshal(userTest)
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/foods/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}

		testFoodHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Update Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataFood()[0])
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/foods/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("UpdateFood", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		testFoodHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

}

func TestFoodRouter_DeleteHandler(t *testing.T) {

	t.Run("Error Param Delete Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/api/v1/foods/{id}", nil)

		mockRepository := &repoMock.FoodRepository{}
		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}

		testFoodHandler.DeleteHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Delete Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodDelete, "/api/v1/foods/{id}", nil)
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("DeleteFood", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error sql")).Once()

		testFoodHandler.DeleteHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Delete Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodDelete, "/api/v1/foods/{id}", nil)
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.FoodRepository{}

		testFoodHandler := &v1.FoodRouter{Repo: mockRepository}
		mockRepository.On("DeleteFood", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		testFoodHandler.DeleteHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

}
