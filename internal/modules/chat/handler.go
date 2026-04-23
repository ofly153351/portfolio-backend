package chat

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	ws "github.com/gofiber/websocket/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/chat/ws", ws.New(h.ChatWS))
}

func (h *Handler) ChatWS(conn *ws.Conn) {
	send := func(evt StreamEvent) error {
		return conn.WriteJSON(evt)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var req ChatRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			_ = send(StreamEvent{Type: "error", Error: "invalid request body"})
			continue
		}

		_ = send(StreamEvent{Type: "status", Message: "streaming"})
		err = h.service.StreamAsk(context.Background(), req, func(evt StreamEvent) error {
			return send(evt)
		})
		if err != nil {
			_ = send(StreamEvent{Type: "error", Error: err.Error()})
			continue
		}
	}
}
