package integration

import (
	"food-api/infrastructure"
	"food-api/test/integration/seed"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"

	data "food-api/infrastructure/database"
	_ "github.com/joho/godotenv/autoload"
)

// a is a reference to the main Application type. This is used for its database
// connection that it harbours inside of the type as well as the route definitions
// that are defined on the embedded handler.
var server *infrastructure.Server
var dataConnection *data.Data


// TestMain calls testMain and passes the returned exit code to os.Exit(). The reason
// that TestMain is basically a wrapper around testMain is because os.Exit() does not
// respect deferred functions, so this configuration allows for a deferred function.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// testMain returns an integer denoting an exit code to be returned and used in
// TestMain. The exit code 0 denotes success, all other codes denote failure (1
// and 2).
func testMain(m *testing.M) int {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	//Connect data base
	dbc := seed.Open()
	defer data.CloseTest()

	//Versioning the database
	err = data.VersionedDB(dbc, true)
	if err != nil {
		log.Fatal(err)
	}

	//redis details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redis, err := data.NewRedisDB(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	// Run Server for test
	port := os.Getenv("DAEMON_PORT")
	server = infrastructure.NewServerTest(port, dbc, redis)
	dataConnection = dbc


	return m.Run()
}
