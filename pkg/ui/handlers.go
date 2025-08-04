package ui

import (
	"log"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func SetupUIHandlers(app *fiber.App) {
	app.Get("/", templAdaptor(dashboard()))

	app.Get("/study/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		return templAdaptor(studySummaryPage(uint(id)))(c)
	})

	app.Get("/study/:id/history", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		return templAdaptor(studyHistory(uint(id)))(c)
	})

	app.Get("/study/:id/trials", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		return templAdaptor(trialsListPage(uint(id)))(c)
	})
}

func templAdaptor(component templ.Component) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		err := component.Render(c.Context(), c.Response().BodyWriter())
		if err != nil {
			log.Print(err)
		}

		return err
	}
}
