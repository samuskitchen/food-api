package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"food-api/domain/user/application/v1/response"
	"food-api/domain/user/domain/model"
	repoDomain "food-api/domain/user/domain/repository"
	"food-api/domain/user/infrastructure/persistence"
	"food-api/infrastructure/database"
	"food-api/infrastructure/middleware"
	"github.com/go-chi/chi"
	"net/http"
	"time"
)

// UserRouter
type UserRouter struct {
	Repo repoDomain.UserRepository
}

// NewUserHandler
func NewUserHandler(db *database.Data) *UserRouter {
	return &UserRouter{
		Repo: persistence.NewUserRepository(db),
	}
}

// GetAllUser response all the users.
func (ur *UserRouter) GetAllUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := ur.Repo.GetAllUser(ctx)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusOK, users)
}

// GetOneHandler response one user by id.
func (ur *UserRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	ctx := r.Context()
	userResult, err := ur.Repo.GetById(ctx, id)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusOK, userResult)
}

// CreateHandler Create a new user.
func (ur *UserRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	defer r.Body.Close()
	if err := user.HashPassword(); err != nil {
		_ = middleware.HTTPError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	userErrors := user.Validate("")
	if len(userErrors) > 0 {
		_ = middleware.HTTPErrors(w, r, http.StatusUnprocessableEntity, userErrors)
		return
	}


	ctx := r.Context()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := ur.Repo.CreateUser(ctx, &user)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s%s", r.URL.String(), result.ID))
	_ = middleware.JSON(w, r, http.StatusCreated, result)
}

// UpdateHandler update a stored user by id.
func (ur *UserRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	id := chi.URLParam(r, "id")

	if id == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	var userUpdate model.User
	err := json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()
	userErrors := userUpdate.Validate("update")
	if len(userErrors) > 0 {
		_ = middleware.HTTPErrors(w, r, http.StatusUnprocessableEntity, userErrors)
		return
	}

	ctx := r.Context()
	userUpdate.UpdatedAt = now

	err = ur.Repo.UpdateUser(ctx, id, userUpdate)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	result := response.UserResponse{
		ID:        userUpdate.ID,
		Names:     userUpdate.Names,
		LastNames: userUpdate.LastNames,
		Email:     userUpdate.Email,
	}

	_ = middleware.JSON(w, r, http.StatusOK, result)
}
