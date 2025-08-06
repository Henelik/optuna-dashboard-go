package main

import (
	"log"
	"os"

	"github.com/Henelik/optuna-dashboard-go/pkg/db"
	"github.com/Henelik/optuna-dashboard-go/pkg/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %s", err.Error())
	}

	database, err := gorm.Open(postgres.Open(os.Getenv("OPTUNA_DB")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}

	db.DB = database

	app := fiber.New()
	app.Use(logger.New())
	ui.SetupUIHandlers(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
