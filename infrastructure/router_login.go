package infrastructure

import (
	userApp "food-api/domain/user/application"
	"food-api/infrastructure/auth"
	"food-api/infrastructure/database"
	"github.com/go-chi/chi"
	"net/http"
)

// Routes returns the API Handler with configuration.
func RoutesLogin(conn *database.Data, redis *database.RedisService) http.Handler {
	router := chi.NewRouter()

	tk := auth.NewToken()
	lr := userApp.NewLoginHandler(conn, redis, tk)
	router.Mount("/", routesLogin(lr))

	return router
}

// routesFood returns login router with each endpoint.
func routesLogin(handler *userApp.LoginRouter) http.Handler {
	router := chi.NewRouter()

	router.Post("/login", handler.LoginHandler)
	router.Post("/logout", handler.LogoutHandler)
	router.Post("/refresh", handler.RefreshHandler)

	return router
}