package main

import (
	"log"
	"os"

	"github.com/Henelik/optuna-dashboard-go/pkg/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("OPTUNA_DB")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}

	ui.DB = db

	app := fiber.New()

	ui.SetupUIHandlers(app)

	log.Fatal(app.Listen(":3000"))
}
