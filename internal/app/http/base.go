package http

import (
	"github.com/gofiber/fiber/v2"
	"portfolio-backend/internal/config"
)

func registerBaseRoutes(app *fiber.App, cfg config.Config) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "AI Portfolio Backend is running",
			"app":     cfg.AppName,
		})
	})
}
