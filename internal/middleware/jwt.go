package middleware

import "github.com/gofiber/fiber/v2"

// RequireJWT is a placeholder request guard for protected routes.
func RequireJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
