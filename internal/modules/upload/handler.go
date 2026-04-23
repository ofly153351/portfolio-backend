package upload

import (
	"errors"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	authmodule "portfolio-backend/internal/modules/auth"
)

type Handler struct {
	service *Service
	auth    *authmodule.Handler
}

func NewHandler(service *Service, auth *authmodule.Handler) *Handler {
	return &Handler{service: service, auth: auth}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/admin/upload", h.UploadImage)
	router.Post("/admin/technical/upload", h.UploadTechnicalAsset)
}

func (h *Handler) UploadImage(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "upload_service_unavailable"})
	}
	if _, err := h.auth.ValidateRequest(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	form, _ := c.MultipartForm()
	files := make([]*multipart.FileHeader, 0)
	if form != nil {
		if list, ok := form.File["files"]; ok {
			files = append(files, list...)
		}
	}
	if single, err := c.FormFile("file"); err == nil && single != nil {
		files = append(files, single)
	}
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file_required"})
	}

	urls, err := h.uploadFiles(c, files, "projects")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(UploadImageResponse{
		OK:   true,
		URL:  urls[0],
		URLs: urls,
	})
}

func (h *Handler) UploadTechnicalAsset(c *fiber.Ctx) error {
	if h.service == nil || h.auth == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "upload_service_unavailable"})
	}
	if _, err := h.auth.ValidateRequest(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	form, _ := c.MultipartForm()
	files := make([]*multipart.FileHeader, 0)
	if form != nil {
		if list, ok := form.File["files"]; ok {
			files = append(files, list...)
		}
	}
	if single, err := c.FormFile("file"); err == nil && single != nil {
		files = append(files, single)
	}
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file_required"})
	}

	urls, err := h.uploadFiles(c, files, "technical")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(UploadImageResponse{
		OK:   true,
		URL:  urls[0],
		URLs: urls,
	})
}

func (h *Handler) uploadFiles(c *fiber.Ctx, files []*multipart.FileHeader, folder string) ([]string, error) {
	urls := make([]string, 0, len(files))
	for _, file := range files {
		url, err := h.service.UploadImage(c.UserContext(), file, folder)
		if err != nil {
			if errors.Is(err, ErrFileRequired) {
				return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file_required"})
			}
			if errors.Is(err, ErrInvalidFileType) {
				return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_file_type"})
			}
			return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "upload_failed"})
		}
		urls = append(urls, url)
	}
	return urls, nil
}
