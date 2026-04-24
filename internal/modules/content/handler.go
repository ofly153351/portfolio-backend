package content

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	authmodule "portfolio-backend/internal/modules/auth"
	publicauthmodule "portfolio-backend/internal/modules/publicauth"
)

type Handler struct {
	service    *Service
	auth       *authmodule.Handler
	publicAuth *publicauthmodule.Handler
}

func NewHandler(service *Service, auth *authmodule.Handler, publicAuth *publicauthmodule.Handler) *Handler {
	return &Handler{service: service, auth: auth, publicAuth: publicAuth}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/admin/content", h.GetAdminContent)
	router.Put("/admin/content", h.PutAdminContent)
	router.Post("/admin/content/publish", h.PublishContent)
	router.Get("/admin/content/history", h.GetContentHistory)
	router.Get("/admin/technical", h.GetTechnical)
	router.Post("/admin/technical", h.CreateTechnical)
	router.Put("/admin/technical/:id", h.UpdateTechnical)
	router.Delete("/admin/technical/:id", h.DeleteTechnical)
	router.Get("/content", h.GetPublicContent)
}

func (h *Handler) GetAdminContent(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	if _, err := h.auth.ValidateRequest(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}

	locale := c.Query("locale", "en")
	resp, err := h.service.GetDraft(locale)
	if err != nil {
		return mapContentError(c, err, 0)
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) PutAdminContent(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	actor, err := h.auth.ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}

	var req PutAdminContentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_request_body"})
	}

	locale := c.Query("locale", "en")
	resp, err := h.service.SaveDraft(locale, req, actor)
	if err != nil {
		currentVersion := 0
		if errors.Is(err, ErrVersionConflict) {
			current, currentErr := h.service.GetDraft(locale)
			if currentErr == nil {
				currentVersion = current.Version
			}
		}
		return mapContentError(c, err, currentVersion)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) PublishContent(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	actor, err := h.auth.ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	locale := c.Query("locale", "en")
	resp, err := h.service.Publish(locale, actor)
	if err != nil {
		return mapContentError(c, err, 0)
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) GetContentHistory(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	if _, err := h.auth.ValidateRequest(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	locale := c.Query("locale", "en")
	resp, err := h.service.History(locale)
	if err != nil {
		return mapContentError(c, err, 0)
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) GetPublicContent(c *fiber.Ctx) error {
	if h.service == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	if h.publicAuth != nil {
		if err := h.publicAuth.ValidateHTTP(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "public_token_invalid"})
		}
	}
	locale := c.Query("locale", "en")
	resp, err := h.service.GetPublished(locale)
	if err != nil {
		return mapContentError(c, err, 0)
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) GetTechnical(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	if _, err := h.auth.ValidateRequest(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	locale := c.Query("locale", "en")
	resp, err := h.service.GetTechnical(locale)
	if err != nil {
		return mapContentError(c, err, 0)
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) CreateTechnical(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	actor, err := h.auth.ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	var item TechnicalItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_request_body"})
	}
	locale := c.Query("locale", "en")
	resp, err := h.service.CreateTechnical(locale, item, actor)
	if err != nil {
		return mapContentError(c, err, 0)
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) UpdateTechnical(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	actor, err := h.auth.ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	id := c.Params("id")
	var item TechnicalItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_request_body"})
	}
	locale := c.Query("locale", "en")
	resp, err := h.service.UpdateTechnical(locale, id, item, actor)
	if err != nil {
		return mapContentError(c, err, 0)
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) DeleteTechnical(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "content_service_unavailable"})
	}
	actor, err := h.auth.ValidateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	id := c.Params("id")
	locale := c.Query("locale", "en")
	resp, err := h.service.DeleteTechnical(locale, id, actor)
	if err != nil {
		return mapContentError(c, err, 0)
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func mapContentError(c *fiber.Ctx, err error, currentVersion int) error {
	switch {
	case errors.Is(err, ErrInvalidLocale):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_locale"})
	case errors.Is(err, ErrVersionConflict):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":           "version_conflict",
			"current_version": currentVersion,
		})
	case errors.Is(err, ErrInvalidPayload):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_payload"})
	case errors.Is(err, ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not_found"})
	case errors.Is(err, ErrUnauthorized):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal_error"})
	}
}
