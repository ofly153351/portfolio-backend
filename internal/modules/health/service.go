package health

import "context"

// Service contains use cases for the health module.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Health(ctx context.Context) error {
	if s.repo == nil {
		return ErrNotImplemented
	}
	return s.repo.Ping(ctx)
}
