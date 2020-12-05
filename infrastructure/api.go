package infrastructure

import (
	"food-api/infrastructure/database"
	"log"
	"os"
)

func Start(port string) {

	// connection to the database.
	db := database.New()
	defer db.DB.Close()

	//Versioning the database
	err := database.VersionedDB(db, false)
	if err != nil {
		log.Fatal(err)
	}

	//redis details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redis, err := database.NewRedisDB(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	server := newServer(port, db, redis)

	// start the server.
	server.Start()
}
