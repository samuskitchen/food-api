package main

import (
	"food-api/infrastructure"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.Println("stating API cmd")
	port := os.Getenv("API_PORT")
	infrastructure.Start(port)
}