package http

import "github.com/gofiber/fiber/v2"

func registerContentRoutes(api fiber.Router, deps *Dependencies) {
	deps.ContentHandler.RegisterRoutes(api)
}
