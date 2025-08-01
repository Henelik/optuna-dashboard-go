package ui

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func SetupUIHandlers(app *fiber.App) {
	app.Get("/", templAdaptor(dashboard))
}

func templAdaptor(component func() templ.Component) fiber.Handler {
	return adaptor.HTTPHandler(templ.Handler(component()))
}
