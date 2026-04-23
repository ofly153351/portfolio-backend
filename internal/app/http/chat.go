package http

import "github.com/gofiber/fiber/v2"

func registerChatRoutes(api fiber.Router, deps *Dependencies) {
	deps.ChatHandler.RegisterRoutes(api)
}
