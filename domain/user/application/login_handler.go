package application

import (
	"encoding/json"
	"fmt"
	"food-api/domain/user/domain/model"
	repoDomain "food-api/domain/user/domain/respository"
	"food-api/domain/user/infrastructure/persistence"
	"food-api/infrastructure/auth"
	authModel "food-api/infrastructure/auth/model"
	"food-api/infrastructure/database"
	"food-api/infrastructure/middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// LoginRouter
type LoginRouter struct {
	Repo  repoDomain.UserRepository
	Redis *database.RedisService
	Token auth.TokenInterface
}

// NewLoginHandler
func NewLoginHandler(db *database.Data, redis *database.RedisService, token auth.TokenInterface) *LoginRouter {
	return &LoginRouter{
		Repo:  persistence.NewUserRepository(db),
		Redis: redis,
		Token: token,
	}
}

// LoginHandler
func (lr *LoginRouter) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	defer r.Body.Close()
	userErrors := user.Validate("login")
	if len(userErrors) > 0 {
		_ = middleware.HTTPErrors(w, r, http.StatusInternalServerError, userErrors)
		return
	}

	ctx := r.Context()
	result, err := lr.Repo.GetUserByEmailAndPassword(ctx, &user)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	tokenDetails, err := lr.Token.CreateToken(result.ID)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := lr.Redis.Auth.CreateAuth(ctx, result.ID, tokenDetails)
	if saveErr != nil {
		_ = middleware.HTTPErrors(w, r, http.StatusInternalServerError, userErrors)
		return
	}

	userData := authModel.DataLogin{
		ID:           result.ID,
		Names:        result.Names,
		LastNames:    result.LastNames,
		AccessToken:  tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshToken,
	}

	_ = middleware.JSON(w, r, http.StatusOK, userData)
}

// LogoutHandler
func (lr *LoginRouter) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//check is the user is authenticated first
	metadata, err := lr.Token.ExtractTokenMetadata(r)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	//if the access token exist and it is still valid, then delete both the access token and the refresh token
	err = lr.Redis.Auth.DeleteTokens(ctx, metadata)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusOK, "Successfully logged out")
}

// RefreshHandler is the function that uses the refresh_token to generate new pairs of refresh and access tokens.
func (lr *LoginRouter) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var dataLogin authModel.DataLogin
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&dataLogin)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	//verify the refresh Token and is refresh Token valid?
	refreshToken, err := auth.VerifyAndValidateRefreshToken(dataLogin.RefreshToken)
	if err != nil {
		//any error may be due to refreshToken expiration
		_ = middleware.HTTPError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if ok && refreshToken.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string

		if !ok {
			_ = middleware.HTTPError(w, r, http.StatusUnprocessableEntity, "Cannot get uuid")
			return
		}

		userId := fmt.Sprintf("%.f", claims["user_id"])

		//Delete the previous RefreshHandler Token
		delErr := lr.Redis.Auth.DeleteRefresh(ctx, refreshUuid)
		if delErr != nil {
			//if any goes wrong
			_ = middleware.HTTPError(w, r, http.StatusUnauthorized, "unauthorized")
			return
		}

		//Create new pairs of refresh and access tokens
		tokenDetails, err := lr.Token.CreateToken(userId)
		if err != nil {
			_ = middleware.HTTPError(w, r, http.StatusForbidden, err.Error())
			return
		}

		//save the tokens metadata to redis
		err = lr.Redis.Auth.CreateAuth(ctx, userId, tokenDetails)
		if err != nil {
			_ = middleware.HTTPError(w, r, http.StatusForbidden, err.Error())
			return
		}

		dataLogin.RefreshToken = tokenDetails.RefreshToken
		dataLogin.AccessToken = tokenDetails.AccessToken

		_ = middleware.JSON(w, r, http.StatusOK, dataLogin)
	} else {
		_ = middleware.HTTPError(w, r, http.StatusUnauthorized, "refresh token expired")
	}
}
