package http

import "github.com/gofiber/fiber/v2"

func registerPublicRoutes(api fiber.Router, deps *Dependencies) {
	deps.PublicHandler.RegisterRoutes(api)
}

