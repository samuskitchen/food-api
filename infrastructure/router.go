package infrastructure

import (
	//"food-api/domain/user/application/v1"
	"food-api/infrastructure/database"
	"github.com/go-chi/chi"
	"net/http"
)

// Routes returns the API V1 Handler with configuration.
func Routes(conn *database.Data) http.Handler {
	router := chi.NewRouter()

	//ur := v1.NewSatellitesHandler(conn)
	//router.Mount("/", routesSatellite(ur))

	return router
}

// routesUser returns user router with each endpoint.
/*func routesUser(handler *v1.SatellitesRouter) http.Handler {
	router := chi.NewRouter()

	router.Get("/topsecret_split", handler.TopSecretSplitHandler)
	router.Post("/topsecret", handler.TopSecretHandler)
	router.Post("/topsecret_split/{satellite_name}", handler.TopSecretSplitSatelliteNameHandler)

	return router
}*/
