// Package main Food API.
//
// The purpose of this application is to provide food service by users
//
//
// This should show the struct of endpoints
// Terms Of Service:
//
//     Schemes: http, https
//     Host: localhost:8888
//     BasePath: /api/v1
//     Version: 1.0.0
//     Contact: https://www.linkedin.com/in/daniel-de-la-pava-suarez/
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
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