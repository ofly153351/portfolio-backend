package http

import "github.com/gofiber/fiber/v2"

func registerHealthRoutes(api fiber.Router, deps *Dependencies) {
	deps.HealthHandler.RegisterRoutes(api)
}
