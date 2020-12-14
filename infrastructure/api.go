package infrastructure

import (
	"food-api/infrastructure/database"
	"log"
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
	redis := database.NewRedisDB()

	server := newServer(port, db, redis)

	// start the server.
	server.Start()
}
