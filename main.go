package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/samuel-fonseca/task-manager-api/database"
	"github.com/samuel-fonseca/task-manager-api/router"
)

func main() {
	// load the env variables
	loadEnvVariables()
	// connect the database
	database.ConnectDatabase()
	// boot server
	startServer()
}

func loadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func startServer() *fiber.App {
	port := os.Getenv("CLIENT_PORT")

	app := fiber.New()
	router.RegisterRoutes(app)
	app.Listen(":" + port)

	return app
}
