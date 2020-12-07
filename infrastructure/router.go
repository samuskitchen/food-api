package infrastructure

import (
	v1Food "food-api/domain/food/application/v1"
	v1User "food-api/domain/user/application/v1"
	"food-api/infrastructure/database"
	"food-api/infrastructure/middleware"
	"github.com/go-chi/chi"
	"net/http"
)

// Routes returns the API V1 Handler with configuration.
func Routes(conn *database.Data) http.Handler {
	router := chi.NewRouter()

	ur := v1User.NewUserHandler(conn)
	router.Mount("/users", routesUser(ur))

	fr := v1Food.NewFoodHandler(conn)
	router.With(middleware.AuthMiddleware).Mount("/foods", routesFood(fr))

	return router
}

// routesUser returns user router with each endpoint.
func routesUser(handler *v1User.UserRouter) http.Handler {
	router := chi.NewRouter()

	router.With(middleware.AuthMiddleware).Get("/", handler.GetAllUserHandler)
	router.With(middleware.AuthMiddleware).Get("/{id}", handler.GetOneHandler)
	router.Post("/", handler.CreateHandler)
	router.With(middleware.AuthMiddleware).Put("/{id}", handler.UpdateHandler)

	return router
}

// routesFood returns food router with each endpoint.
func routesFood(handler *v1Food.FoodRouter) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handler.GetAllFoodHandler)
	router.Get("/{id}", handler.GetOneHandler)
	router.Get("/user/{id}", handler.GetOneByUserHandler)
	router.With(middleware.MaxSizeAllowed).Post("/", handler.CreateHandler)
	router.With(middleware.MaxSizeAllowed).Put("/{id}", handler.UpdateHandler)
	router.Delete("/{id}", handler.DeleteHandler)

	return router
}
