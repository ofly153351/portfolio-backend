package auth

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
	mu      sync.Mutex
	hits    map[string][]time.Time
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
		hits:    make(map[string][]time.Time),
	}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/admin/login", h.Login)
	router.Post("/admin/logout", h.Logout)
	router.Get("/admin/me", h.Me)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	if h.service == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth_service_unavailable"})
	}
	if !h.allowLogin(c.IP()) {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "rate_limited"})
	}
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_request_body"})
	}

	token, user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		if err == ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "login_failed"})
	}

	c.Set("Authorization", "Bearer "+token)
	return c.Status(fiber.StatusOK).JSON(LoginResponse{
		OK:          true,
		AccessToken: token,
		TokenType:   "Bearer",
		User:        user,
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	if h.service == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth_service_unavailable"})
	}
	token := extractBearerToken(c.Get("Authorization"))
	h.service.Logout(token)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ok": true})
}

func (h *Handler) Me(c *fiber.Ctx) error {
	if h.service == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth_service_unavailable"})
	}
	token := extractBearerToken(c.Get("Authorization"))
	user, err := h.service.Validate(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(MeResponse{Authenticated: false})
	}
	return c.Status(fiber.StatusOK).JSON(MeResponse{Authenticated: true, User: &user})
}

func (h *Handler) ValidateRequest(c *fiber.Ctx) (User, error) {
	if h.service == nil {
		return User{}, ErrUnauthorized
	}
	token := extractBearerToken(c.Get("Authorization"))
	return h.service.Validate(token)
}

func (h *Handler) allowLogin(ip string) bool {
	if h == nil {
		return false
	}
	now := time.Now()
	window := now.Add(-1 * time.Minute)
	const limit = 10

	h.mu.Lock()
	defer h.mu.Unlock()

	series := h.hits[ip]
	filtered := make([]time.Time, 0, len(series)+1)
	for _, t := range series {
		if t.After(window) {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) >= limit {
		h.hits[ip] = filtered
		return false
	}
	filtered = append(filtered, now)
	h.hits[ip] = filtered
	return true
}
