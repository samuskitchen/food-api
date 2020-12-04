package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"food-api/domain/food/application/v1/response"
	"food-api/domain/food/domain/model"
	repoDomain "food-api/domain/food/domain/repository"
	"food-api/domain/food/infrastructure/persistence"
	"food-api/infrastructure/database"
	"food-api/infrastructure/middleware"
	"github.com/go-chi/chi"
	"net/http"
	"time"
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

// GetAllFood response all the food.
func (ur *FoodRouter) GetAllFood(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := ur.Repo.GetAllFood(ctx)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusOK, users)
}

// GetOneHandler response one food by id.
func (ur *FoodRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	ctx := r.Context()
	userResult, err := ur.Repo.GetFoodById(ctx, id)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusOK, userResult)
}

// GetOneByUserHandler response one food by user id.
func (ur *FoodRouter) GetOneByUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	ctx := r.Context()
	userResult, err := ur.Repo.GetFoodByUserId(ctx, id)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusOK, userResult)
}

// CreateHandler Create a new food.
func (ur *FoodRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	var food model.Food

	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()
	foodErrors := food.Validate("")
	if len(foodErrors) > 0 {
		_ = middleware.HTTPErrors(w, r, http.StatusUnprocessableEntity, foodErrors)
		return
	}

	ctx := r.Context()
	food.CreatedAt = now
	food.UpdatedAt = now

	result, err := ur.Repo.SaveFood(ctx, &food)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s%s", r.URL.String(), result.ID))
	_ = middleware.JSON(w, r, http.StatusCreated, result)
}

// UpdateHandler update a stored food by id.
func (ur *FoodRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	var foodUpdate model.Food
	err := json.NewDecoder(r.Body).Decode(&foodUpdate)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()
	userErrors := foodUpdate.Validate("update")
	if len(userErrors) > 0 {
		_ = middleware.HTTPErrors(w, r, http.StatusUnprocessableEntity, userErrors)
		return
	}

	ctx := r.Context()
	foodUpdate.UpdatedAt = now

	err = ur.Repo.UpdateFood(ctx, id, &foodUpdate)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	result := response.FoodResponse{
		ID:          foodUpdate.ID,
		UserID:      foodUpdate.UserID,
		Title:       foodUpdate.Title,
		Description: foodUpdate.Description,
		FoodImage:   foodUpdate.FoodImage,
	}

	_ = middleware.JSON(w, r, http.StatusOK, result)
}

// DeleteHandler Remove a food by ID.
func (ur *FoodRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	ctx := r.Context()
	err := ur.Repo.DeleteFood(ctx, id)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusNoContent, middleware.Map{})
}