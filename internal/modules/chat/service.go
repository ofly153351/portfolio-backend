package chat

import "context"

// Service contains use cases for the chat module.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Ask(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	if req.Message == "" {
		return ChatResponse{}, ErrMessageRequired
	}
	if s.repo == nil {
		return ChatResponse{}, ErrNotImplemented
	}
	return s.repo.Chat(ctx, req)
}

func (s *Service) StreamAsk(ctx context.Context, req ChatRequest, onEvent func(StreamEvent) error) error {
	if req.Message == "" {
		return ErrMessageRequired
	}
	if s.repo == nil {
		return ErrNotImplemented
	}
	return s.repo.StreamChat(ctx, req, onEvent)
}
