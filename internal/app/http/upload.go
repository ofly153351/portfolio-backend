package http

import "github.com/gofiber/fiber/v2"

func registerUploadRoutes(api fiber.Router, deps *Dependencies) {
	deps.UploadHandler.RegisterRoutes(api)
}
