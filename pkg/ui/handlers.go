package ui

import (
	"log"
	"strconv"

	"github.com/Henelik/optuna-dashboard-go/pkg/db"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

		anchorID, err := strconv.Atoi(c.Query("t", "-1"))
		if err != nil {
			return err
		}

		return templAdaptor(trialsListPage(uint(id), anchorID))(c)
	})

	app.Get("/study/:id/trials/:page", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		page, err := c.ParamsInt("page")
		if err != nil {
			return err
		}

		anchorID, err := strconv.Atoi(c.Query("t", "-1"))
		if err != nil {
			return err
		}

		return templAdaptor(getTrialsRows(uint(id), anchorID, page))(c)
	})

	app.Delete("/study/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		if err := db.DB.Transaction(func(tx *gorm.DB) error {
			return db.DeleteStudy(uint(id), tx)
		}); err != nil {
			return err
		}

		c.Response().Header.Add("hx-redirect", "/")

		return c.SendStatus(fiber.StatusOK)
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
