package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"news-service/config"
	"news-service/database"
	"news-service/handlers"
)

func main() {
	app := fiber.New()

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load config", err)
	}

	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("Failed to connect to db", err)
	}

	app.Post("/edit/:Id", handlers.EditNews)
	app.Get("/list", handlers.ListNews)

	log.Fatal(app.Listen(":8080"))
}
