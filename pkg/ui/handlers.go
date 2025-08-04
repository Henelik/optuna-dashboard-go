package ui

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func SetupUIHandlers(app *fiber.App) {
	app.Get("/", templAdaptor(dashboard))

	app.Get("/study/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		return adaptor.HTTPHandler(templ.Handler(studyHistory(uint(id))))(c)
	})

	app.Get("/study/:id/trials", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		return adaptor.HTTPHandler(templ.Handler(trialsListPage(uint(id))))(c)
	})
}

func templAdaptor(component func() templ.Component) fiber.Handler {
	return adaptor.HTTPHandler(templ.Handler(component()))
}
