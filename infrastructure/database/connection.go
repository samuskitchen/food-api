package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	// registering database driver
	_ "github.com/lib/pq"
)

var (
	data *Data
	once sync.Once
)

type Data struct {
	DB *sql.DB
}

// New returns a new instance of Data with the database connection ready.
func New() *Data {
	once.Do(initDB)
	return data
}

func NewTest() *Data {
	once.Do(initDBTest)
	return data
}

func initDB() {
	db, err := getConnection()
	if err != nil {
		log.Println("Cannot connect to database")
		log.Fatal("This is the error:", err)
	} else {
		log.Println("We are connected to the database")
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	data = &Data{
		DB: db,
	}
}

func initDBTest() {

	db, err := getConnectionTest()
	if err != nil {
		log.Println("Cannot connect to database test")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database test")
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	data = &Data{
		DB: db,
	}
}

func getConnection() (*sql.DB, error) {

	DbHost := os.Getenv("DB_HOST")
	DbDriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	return sql.Open(DbDriver, uri)
}

func getConnectionTest() (*sql.DB, error) {
	DbHost := os.Getenv("DB_HOST")
	DbDriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	return sql.Open(DbDriver, uri)
}

// CloseTest closes the resources used by data test.
func CloseTest() error {
	if data == nil {
		return nil
	}

	return data.DB.Close()
}
