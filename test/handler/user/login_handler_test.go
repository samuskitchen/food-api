package user

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"food-api/domain/user/application"
	responseUser "food-api/domain/user/application/v1/response"
	"food-api/domain/user/domain/model"
	repoMock "food-api/domain/user/domain/repository/mocks"
	authMock "food-api/infrastructure/auth/mocks"
	modelAuth "food-api/infrastructure/auth/model"
	"food-api/infrastructure/database"
	"github.com/alicebob/miniredis"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestRedis returns a redis.Client.
func newTestRedis() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

// generateToken returns a unique token based on the provided email string
func generateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	harsher := md5.New()
	harsher.Write(hash)
	return hex.EncodeToString(harsher.Sum(nil))
}

// dataUser is data for test
func dataUser() model.User {
	now := time.Now()

	return model.User{
		ID:        uuid.New().String(),
		Names:     "Daniel",
		LastNames: "De La Pava Suarez",
		Email:     "daniel.delapava@jikkosoft.com",
		Password:  "123456",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func dataLogin() modelAuth.DataLogin {
	return modelAuth.DataLogin{
		RefreshToken: generateToken("daniel.delapava@jikkosoft.com"),
		AccessToken: generateToken("daniel.delapava@jikkosoft.com"),
	}
}

func dataAccessDetails() *modelAuth.AccessDetails {
	return &modelAuth.AccessDetails{
		UserId: uuid.New().String(),
		TokenUuid: uuid.New().String(),
	}
}

func dataTokenDetails() *modelAuth.TokenDetails {
	return &modelAuth.TokenDetails{
		AccessToken:  generateToken("daniel.delapava@jikkosoft.com"),
		RefreshToken: generateToken("daniel.delapava@jikkosoft.com"),
		TokenUuid:    uuid.New().String(),
		RefreshUuid:  uuid.New().String(),
		AtExpires:    time.Now().Add(time.Minute * 15).Unix(),
		RtExpires:    time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
}

func dataJwt() *jwt.Token {
	claimsData := jwt.MapClaims{}
	claimsData["user_id"] = dataLogin().ID
	claimsData["refresh_uuid"] = dataTokenDetails().RefreshUuid

	return &jwt.Token{
		Raw:       "",
		Method:    nil,
		Header:    nil,
		Claims:    claimsData,
		Signature: "",
		Valid:     true,
	}
}

func dataJwtErrorRefresh() *jwt.Token {
	claimsData := jwt.MapClaims{}
	claimsData["user_id"] = dataLogin().ID

	return &jwt.Token{
		Raw:       "",
		Method:    nil,
		Header:    nil,
		Claims:    claimsData,
		Signature: "",
		Valid:     true,
	}
}

func dataJwtErrorUserId() *jwt.Token {
	claimsData := jwt.MapClaims{}
	claimsData["refresh_uuid"] = dataTokenDetails().RefreshUuid

	return &jwt.Token{
		Raw:       "",
		Method:    nil,
		Header:    nil,
		Claims:    claimsData,
		Signature: "",
		Valid:     true,
	}
}

// dataUserResponse is data for test
func dataUserResponse() *responseUser.UserResponse {

	return &responseUser.UserResponse{
		ID:        uuid.New().String(),
		Names:     "Daniel",
		LastNames: "De La Pava Suarez",
		Email:     "daniel.delapava@jikkosoft.com",
	}
}

func TestLoginRouter_LoginHandler(t *testing.T) {

	t.Run("Error Body Login Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(nil))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository}

		testLoginHandler.LoginHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Validate Login Handler", func(tt *testing.T) {

		var userTest = dataUser()
		userTest.Names = ""
		userTest.LastNames = ""
		userTest.Password = ""
		userTest.Email = ""

		marshal, err := json.Marshal(userTest)
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository}

		testLoginHandler.LoginHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Get User By Email And Password Login Handler", func(tt *testing.T) {
		marshal, err := json.Marshal(dataUser())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository}
		mockRepository.On("GetUserByEmailAndPassword", mock.Anything, mock.Anything).Return(nil, errors.New("error sql"))

		testLoginHandler.LoginHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Create Token", func(tt *testing.T) {
		marshal, err := json.Marshal(dataUser())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockAuth := &authMock.InterfaceAuth{}
		mockToken := &authMock.TokenInterface{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockRepository.On("GetUserByEmailAndPassword", mock.Anything, mock.Anything).Return(dataUserResponse(), nil)
		mockToken.On("CreateToken", mock.Anything).Return(nil, errors.New("error create token"))

		testLoginHandler.LoginHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Create Auth", func(tt *testing.T) {
		marshal, err := json.Marshal(dataUser())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockAuth := &authMock.InterfaceAuth{}
		mockToken := &authMock.TokenInterface{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockRepository.On("GetUserByEmailAndPassword", mock.Anything, mock.Anything).Return(dataUserResponse(), nil)
		mockToken.On("CreateToken", mock.Anything).Return(dataTokenDetails(), nil)
		mockAuth.On("CreateAuth", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error create auth"))

		testLoginHandler.LoginHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Login Successfully", func(tt *testing.T) {
		marshal, err := json.Marshal(dataUser())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockAuth := &authMock.InterfaceAuth{}
		mockToken := &authMock.TokenInterface{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockRepository.On("GetUserByEmailAndPassword", mock.Anything, mock.Anything).Return(dataUserResponse(), nil)
		mockToken.On("CreateToken", mock.Anything).Return(dataTokenDetails(), nil)
		mockAuth.On("CreateAuth", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		testLoginHandler.LoginHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})
}

func TestLoginRouter_LogoutHandler(t *testing.T) {

	t.Run("Error Extract Token Metadata", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewReader(nil))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Token: mockToken}
		mockToken.On("ExtractTokenMetadata", mock.Anything).Return(nil, errors.New("error token valid"))

		testLoginHandler.LogoutHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Delete Token", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewReader(nil))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("ExtractTokenMetadata", mock.Anything).Return(dataAccessDetails(), nil)
		mockAuth.On("DeleteTokens", mock.Anything, mock.Anything ).Return(errors.New("error delete token"))

		testLoginHandler.LogoutHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Logout Successfully", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewReader(nil))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("ExtractTokenMetadata", mock.Anything).Return(dataAccessDetails(), nil)
		mockAuth.On("DeleteTokens", mock.Anything, mock.Anything ).Return(nil)

		testLoginHandler.LogoutHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

}

func TestLoginRouter_RefreshHandler(t *testing.T) {

	t.Run("Error Body Refresh Handler", func(tt *testing.T) {

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(nil))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository}

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Verify And Validate Refresh Token", func(tt *testing.T) {
		marshal, err := json.Marshal(dataLogin())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("VerifyAndValidateRefreshToken", mock.Anything).Return(nil, errors.New("error validate and verify token"))

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Token Expired", func(tt *testing.T) {
		marshal, err := json.Marshal(dataLogin())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("VerifyAndValidateRefreshToken", mock.Anything).Return(&jwt.Token{}, nil)

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Refresh uuid", func(tt *testing.T) {
		marshal, err := json.Marshal(dataLogin())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("VerifyAndValidateRefreshToken", mock.Anything).Return(dataJwtErrorRefresh(), nil)

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error User uuid", func(tt *testing.T) {
		marshal, err := json.Marshal(dataLogin())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("VerifyAndValidateRefreshToken", mock.Anything).Return(dataJwtErrorUserId(), nil)

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Delete Refresh", func(tt *testing.T) {
		marshal, err := json.Marshal(dataLogin())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("VerifyAndValidateRefreshToken", mock.Anything).Return(dataJwt(), nil)
		mockAuth.On("DeleteRefresh", mock.Anything, mock.Anything).Return(errors.New("error delete refresh token"))

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Create Token", func(tt *testing.T) {
		marshal, err := json.Marshal(dataLogin())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("VerifyAndValidateRefreshToken", mock.Anything).Return(dataJwt(), nil)
		mockAuth.On("DeleteRefresh", mock.Anything, mock.Anything).Return(nil)
		mockToken.On("CreateToken", mock.Anything).Return(nil, errors.New("error create token"))

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Create Auth", func(tt *testing.T) {
		marshal, err := json.Marshal(dataLogin())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("VerifyAndValidateRefreshToken", mock.Anything).Return(dataJwt(), nil)
		mockAuth.On("DeleteRefresh", mock.Anything, mock.Anything).Return(nil)
		mockToken.On("CreateToken", mock.Anything).Return(dataTokenDetails(), nil)
		mockAuth.On("CreateAuth", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error create auth"))

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Refresh Successfully", func(tt *testing.T) {
		marshal, err := json.Marshal(dataLogin())
		assert.NoError(tt, err)

		request := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(marshal))
		response := httptest.NewRecorder()
		mockRepository := &repoMock.UserRepository{}
		mockToken := &authMock.TokenInterface{}
		mockAuth := &authMock.InterfaceAuth{}

		mockRedis := &database.RedisService{
			Client: newTestRedis(),
			Auth:   mockAuth,
		}

		testLoginHandler := &application.LoginRouter{Repo: mockRepository, Redis: mockRedis, Token: mockToken}
		mockToken.On("VerifyAndValidateRefreshToken", mock.Anything).Return(dataJwt(), nil)
		mockAuth.On("DeleteRefresh", mock.Anything, mock.Anything).Return(nil)
		mockToken.On("CreateToken", mock.Anything).Return(dataTokenDetails(), nil)
		mockAuth.On("CreateAuth", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		testLoginHandler.RefreshHandler(response, request)
		mockRepository.AssertExpectations(tt)
	})
}
