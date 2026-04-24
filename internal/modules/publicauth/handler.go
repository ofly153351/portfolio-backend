package publicauth

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/public/token", h.IssueToken)
}

func (h *Handler) IssueToken(c *fiber.Ctx) error {
	if h.service == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "public_token_service_unavailable"})
	}
	token, expiresAt, err := h.service.Issue()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "public_token_issue_failed"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":      token,
		"token_type": "Bearer",
		"expires_at": expiresAt,
		"expires_in": int(expiresAt.Sub(time.Now().UTC()).Seconds()),
	})
}

func (h *Handler) ValidateHTTP(c *fiber.Ctx) error {
	if h.service == nil {
		return errors.New("public token service unavailable")
	}
	token := strings.TrimSpace(c.Get("X-Public-Token"))
	if token == "" {
		return ErrTokenMissing
	}
	return h.service.Validate(token)
}

func (h *Handler) ValidateRawToken(token string) error {
	if h.service == nil {
		return errors.New("public token service unavailable")
	}
	normalized := strings.TrimSpace(token)
	if normalized == "" {
		return ErrTokenMissing
	}
	return h.service.Validate(normalized)
}
