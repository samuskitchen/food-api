package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	v1 "food-api/domain/user/application/v1"
	responseUser "food-api/domain/user/application/v1/response"
	"food-api/domain/user/domain/model"
	repoMock "food-api/domain/user/domain/repository/mocks"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// dataUser is data for test
func dataUser() []model.User {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	return []model.User{
		{
			ID:        uuid.New().String(),
			Names:     "Daniel",
			LastNames: "De La Pava Suarez",
			Email:     "daniel.delapava@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uuid.New().String(),
			Names:     "Rebecca",
			LastNames: "Romero",
			Email:     "rebecca.romero@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// dataUserResponse is data for test
func dataUserResponse() []responseUser.UserResponse {

	return []responseUser.UserResponse{
		{
			ID:        uuid.New().String(),
			Names:     "Daniel",
			LastNames: "De La Pava Suarez",
			Email:     "daniel.delapava@jikkosoft.com",
		},
		{
			ID:        uuid.New().String(),
			Names:     "Rebecca",
			LastNames: "Romero",
			Email:     "rebecca.romero@jikkosoft.com",
		},
	}
}

func TestUserRouter_GetAllUser(t *testing.T) {

	t.Run("Error Get All User Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}
		mockRepository.On("GetAllUserHandler", mock.Anything, mock.Anything).Return(nil, errors.New("error trace test"))

		testUserHandler.GetAllUserHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get All User Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}
		mockRepository.On("GetAllUserHandler", mock.Anything).Return(dataUserResponse(), nil)

		testUserHandler.GetAllUserHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})
}

func TestUserRouter_GetOneHandler(t *testing.T) {

	t.Run("Error Param Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/{id}", nil)

		mockRepository := &repoMock.UserRepository{}
		testUserHandler := &v1.UserRouter{Repo: mockRepository}

		testUserHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/{id}", nil)

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}
		mockRepository.On("GetById", mock.Anything, mock.Anything).Return(responseUser.UserResponse{}, errors.New("error sql")).Once()

		testUserHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get One Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/{id}", nil)

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}
		mockRepository.On("GetById", mock.Anything, mock.Anything).Return(responseUser.UserResponse{}, nil).Once()

		testUserHandler.GetOneHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

}

func TestUserRouter_CreateHandler(t *testing.T) {

	t.Run("Error Body Create Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(nil))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}

		testUserHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Create Handler", func(tt *testing.T) {
		dataUser()[0].ID = ""

		marshal, err := json.Marshal(dataUser()[0])
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}
		mockRepository.On("CreateUser", mock.Anything, mock.Anything).Return(&responseUser.UserResponse{}, errors.New("error sql"))

		testUserHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Validate Create Handler", func(tt *testing.T) {

		var userTest = dataUser()[0]
		userTest.Names = ""

		marshal, err := json.Marshal(userTest)
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}

		testUserHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)

	})

	t.Run("Create Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataUser()[0])
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}
		mockRepository.On("CreateUser", mock.Anything, mock.Anything).Return(&dataUserResponse()[0], nil)

		testUserHandler.CreateHandler(response, request)
		mockRepository.AssertExpectations(tt)

	})
}

func TestUserRouter_UpdateHandler(t *testing.T) {

	t.Run("Error Param Update Handler", func(tt *testing.T) {

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", nil)

		mockRepository := &repoMock.UserRepository{}
		testUserHandler := &v1.UserRouter{Repo: mockRepository}

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Body Update Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", bytes.NewReader(nil))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))

		mockRepository := &repoMock.UserRepository{}
		testUserHandler := &v1.UserRouter{Repo: mockRepository}

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Update Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataUser()[0])
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}
		mockRepository.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error sql")).Once()

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Validate Update Handler", func(tt *testing.T) {

		var userTest = dataUser()[0]
		userTest.Email = ""

		marshal, err := json.Marshal(userTest)
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Update Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataUser()[0])
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPut, "/api/v1/users/{id}", bytes.NewReader(marshal))
		response := httptest.NewRecorder()

		requestCtx := chi.NewRouteContext()
		requestCtx.URLParams.Add("id", "1")

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestCtx))
		mockRepository := &repoMock.UserRepository{}

		testUserHandler := &v1.UserRouter{Repo: mockRepository}
		mockRepository.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		testUserHandler.UpdateHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})
}