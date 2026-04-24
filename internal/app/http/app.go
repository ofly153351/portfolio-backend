package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"portfolio-backend/internal/config"
)

func NewApp(cfg config.Config) *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     cfg.CORSAllowMethods,
		AllowHeaders:     cfg.CORSAllowHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
	}))

	registerBaseRoutes(app, cfg)

	deps := NewDependencies(cfg)
	api := app.Group("/api")
	registerAuthRoutes(api, deps)
	registerHealthRoutes(api, deps)
	registerPublicRoutes(api, deps)
	registerContentRoutes(api, deps)
	registerUploadRoutes(api, deps)
	registerChatRoutes(api, deps)

	return app
}
