package main

import (
	"log"

	"github.com/Henelik/optuna-dashboard-go/pkg/ui"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	ui.SetupUIHandlers(app)

	log.Fatal(app.Listen(":3000"))
}
