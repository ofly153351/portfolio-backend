package http

import "github.com/gofiber/fiber/v2"

func registerAuthRoutes(api fiber.Router, deps *Dependencies) {
	deps.AuthHandler.RegisterRoutes(api)
}
