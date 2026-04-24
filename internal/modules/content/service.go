package content

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	authmodule "portfolio-backend/internal/modules/auth"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDraft(locale string) (GetAdminContentResponse, error) {
	normalized := normalizeLocale(locale)
	if !isSupportedLocale(normalized) {
		return GetAdminContentResponse{}, ErrInvalidLocale
	}
	if s.repo == nil {
		return GetAdminContentResponse{}, ErrUnauthorized
	}

	entry, _, err := s.ensureDraft(context.Background(), normalized, authmodule.User{Username: "system"})
	if err != nil {
		return GetAdminContentResponse{}, err
	}
	entry.Content = withContentIndexes(entry.Content)
	return GetAdminContentResponse{
		Locale:    normalized,
		Version:   entry.Version,
		UpdatedAt: entry.UpdatedAt,
		Content:   entry.Content,
	}, nil
}

func (s *Service) SaveDraft(locale string, req PutAdminContentRequest, actor authmodule.User) (PutAdminContentResponse, error) {
	normalized := normalizeLocale(locale)
	if !isSupportedLocale(normalized) {
		return PutAdminContentResponse{}, ErrInvalidLocale
	}
	if s.repo == nil {
		return PutAdminContentResponse{}, ErrUnauthorized
	}
	req.Content = withContentIndexes(req.Content)
	if err := validateContent(req.Content); err != nil {
		return PutAdminContentResponse{}, err
	}

	current, _, err := s.ensureDraft(context.Background(), normalized, actor)
	if err != nil {
		return PutAdminContentResponse{}, err
	}
	if req.Version != current.Version {
		return PutAdminContentResponse{}, ErrVersionConflict
	}

	entry := HistoryItem{
		Locale:    normalized,
		Version:   current.Version + 1,
		Content:   req.Content,
		UpdatedBy: actor.Username,
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.repo.UpsertContent(context.Background(), entry, "draft"); err != nil {
		return PutAdminContentResponse{}, err
	}
	if err := s.repo.AppendHistory(context.Background(), entry); err != nil {
		return PutAdminContentResponse{}, err
	}
	return PutAdminContentResponse{
		OK:        true,
		Version:   entry.Version,
		UpdatedAt: entry.UpdatedAt,
	}, nil
}

func (s *Service) Publish(locale string, actor authmodule.User) (PublishResponse, error) {
	normalized := normalizeLocale(locale)
	if !isSupportedLocale(normalized) {
		return PublishResponse{}, ErrInvalidLocale
	}
	if s.repo == nil {
		return PublishResponse{}, ErrUnauthorized
	}

	current, _, err := s.ensureDraft(context.Background(), normalized, actor)
	if err != nil {
		return PublishResponse{}, err
	}
	published := current
	published.Content = withContentIndexes(published.Content)
	published.UpdatedBy = actor.Username
	published.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpsertContent(context.Background(), published, "published"); err != nil {
		return PublishResponse{}, err
	}
	return PublishResponse{
		OK:               true,
		PublishedVersion: published.Version,
		PublishedAt:      published.UpdatedAt,
	}, nil
}

func (s *Service) GetTechnical(locale string) (TechnicalListResponse, error) {
	admin, err := s.GetDraft(locale)
	if err != nil {
		return TechnicalListResponse{}, err
	}
	return TechnicalListResponse{
		Locale:    admin.Locale,
		Version:   admin.Version,
		UpdatedAt: admin.UpdatedAt,
		Items:     admin.Content.Technical,
	}, nil
}

func (s *Service) CreateTechnical(locale string, item TechnicalItem, actor authmodule.User) (TechnicalMutationResponse, error) {
	normalized := normalizeLocale(locale)
	if !isSupportedLocale(normalized) {
		return TechnicalMutationResponse{}, ErrInvalidLocale
	}
	if s.repo == nil {
		return TechnicalMutationResponse{}, ErrUnauthorized
	}
	if err := validateTechnicalItem(item); err != nil {
		return TechnicalMutationResponse{}, err
	}
	if strings.TrimSpace(item.ID) == "" {
		item.ID = "tech_" + uuid.NewString()
	}

	current, _, err := s.ensureDraft(context.Background(), normalized, actor)
	if err != nil {
		return TechnicalMutationResponse{}, err
	}
	content := current.Content
	content.Technical = append(content.Technical, item)
	content = withContentIndexes(content)
	item = content.Technical[len(content.Technical)-1]

	entry := HistoryItem{
		Locale:    normalized,
		Version:   current.Version + 1,
		Content:   content,
		UpdatedBy: actor.Username,
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.repo.UpsertContent(context.Background(), entry, "draft"); err != nil {
		return TechnicalMutationResponse{}, err
	}
	if err := s.repo.AppendHistory(context.Background(), entry); err != nil {
		return TechnicalMutationResponse{}, err
	}
	return TechnicalMutationResponse{
		OK:        true,
		Version:   entry.Version,
		UpdatedAt: entry.UpdatedAt,
		Item:      item,
	}, nil
}

func (s *Service) UpdateTechnical(locale, id string, item TechnicalItem, actor authmodule.User) (TechnicalMutationResponse, error) {
	normalized := normalizeLocale(locale)
	if !isSupportedLocale(normalized) {
		return TechnicalMutationResponse{}, ErrInvalidLocale
	}
	if s.repo == nil {
		return TechnicalMutationResponse{}, ErrUnauthorized
	}
	targetID := strings.TrimSpace(id)
	if targetID == "" {
		return TechnicalMutationResponse{}, ErrInvalidPayload
	}
	item.ID = targetID
	if err := validateTechnicalItem(item); err != nil {
		return TechnicalMutationResponse{}, err
	}

	current, _, err := s.ensureDraft(context.Background(), normalized, actor)
	if err != nil {
		return TechnicalMutationResponse{}, err
	}
	content := current.Content
	found := false
	for i := range content.Technical {
		if strings.TrimSpace(content.Technical[i].ID) == targetID {
			content.Technical[i] = item
			found = true
			break
		}
	}
	if !found {
		return TechnicalMutationResponse{}, ErrNotFound
	}
	content = withContentIndexes(content)
	for i := range content.Technical {
		if strings.TrimSpace(content.Technical[i].ID) == targetID {
			item = content.Technical[i]
			break
		}
	}

	entry := HistoryItem{
		Locale:    normalized,
		Version:   current.Version + 1,
		Content:   content,
		UpdatedBy: actor.Username,
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.repo.UpsertContent(context.Background(), entry, "draft"); err != nil {
		return TechnicalMutationResponse{}, err
	}
	if err := s.repo.AppendHistory(context.Background(), entry); err != nil {
		return TechnicalMutationResponse{}, err
	}
	return TechnicalMutationResponse{
		OK:        true,
		Version:   entry.Version,
		UpdatedAt: entry.UpdatedAt,
		Item:      item,
	}, nil
}

func (s *Service) DeleteTechnical(locale, id string, actor authmodule.User) (TechnicalDeleteResponse, error) {
	normalized := normalizeLocale(locale)
	if !isSupportedLocale(normalized) {
		return TechnicalDeleteResponse{}, ErrInvalidLocale
	}
	if s.repo == nil {
		return TechnicalDeleteResponse{}, ErrUnauthorized
	}
	targetID := strings.TrimSpace(id)
	if targetID == "" {
		return TechnicalDeleteResponse{}, ErrInvalidPayload
	}

	current, _, err := s.ensureDraft(context.Background(), normalized, actor)
	if err != nil {
		return TechnicalDeleteResponse{}, err
	}
	content := current.Content
	filtered := make([]TechnicalItem, 0, len(content.Technical))
	found := false
	for _, t := range content.Technical {
		if strings.TrimSpace(t.ID) == targetID {
			found = true
			continue
		}
		filtered = append(filtered, t)
	}
	if !found {
		return TechnicalDeleteResponse{}, ErrNotFound
	}
	content.Technical = filtered
	content = withContentIndexes(content)

	entry := HistoryItem{
		Locale:    normalized,
		Version:   current.Version + 1,
		Content:   content,
		UpdatedBy: actor.Username,
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.repo.UpsertContent(context.Background(), entry, "draft"); err != nil {
		return TechnicalDeleteResponse{}, err
	}
	if err := s.repo.AppendHistory(context.Background(), entry); err != nil {
		return TechnicalDeleteResponse{}, err
	}
	return TechnicalDeleteResponse{
		OK:        true,
		Version:   entry.Version,
		UpdatedAt: entry.UpdatedAt,
		DeletedID: targetID,
	}, nil
}

func (s *Service) History(locale string) (HistoryResponse, error) {
	normalized := normalizeLocale(locale)
	if !isSupportedLocale(normalized) {
		return HistoryResponse{}, ErrInvalidLocale
	}
	if s.repo == nil {
		return HistoryResponse{}, ErrUnauthorized
	}
	history, err := s.repo.ListHistory(context.Background(), normalized)
	if err != nil {
		return HistoryResponse{}, err
	}
	for i := range history {
		history[i].Content = withContentIndexes(history[i].Content)
	}
	return HistoryResponse{
		Locale:  normalized,
		History: history,
	}, nil
}

func (s *Service) GetPublished(locale string) (PublicContentResponse, error) {
	normalized := normalizeLocale(locale)
	if !isSupportedLocale(normalized) {
		return PublicContentResponse{}, ErrInvalidLocale
	}
	if s.repo == nil {
		return PublicContentResponse{}, ErrUnauthorized
	}
	entry, found, err := s.repo.GetByLocaleStatus(context.Background(), normalized, "published")
	if err != nil {
		return PublicContentResponse{}, err
	}
	if !found {
		return PublicContentResponse{
			Locale:  normalized,
			Content: defaultContentBody(),
		}, nil
	}
	return PublicContentResponse{
		Locale:  normalized,
		Content: withContentIndexes(entry.Content),
	}, nil
}

func (s *Service) ensureDraft(ctx context.Context, locale string, actor authmodule.User) (HistoryItem, bool, error) {
	entry, found, err := s.repo.GetByLocaleStatus(ctx, locale, "draft")
	if err != nil {
		return HistoryItem{}, false, err
	}
	if found {
		entry.Content = withContentIndexes(entry.Content)
		return entry, true, nil
	}
	now := time.Now().UTC()
	seed := HistoryItem{
		Locale:    locale,
		Version:   1,
		Content:   defaultContentBody(),
		UpdatedBy: actor.Username,
		UpdatedAt: now,
	}
	if strings.TrimSpace(seed.UpdatedBy) == "" {
		seed.UpdatedBy = "system"
	}
	if err := s.repo.UpsertContent(ctx, seed, "draft"); err != nil {
		return HistoryItem{}, false, err
	}
	if err := s.repo.UpsertContent(ctx, seed, "published"); err != nil {
		return HistoryItem{}, false, err
	}
	return seed, false, nil
}

func validateContent(content ContentBody) error {
	if len(content.PortfolioInfo.About) > 5000 {
		return ErrInvalidPayload
	}
	if !isValidURL(content.PortfolioInfo.Github) {
		return ErrInvalidPayload
	}
	if !isValidURL(content.PortfolioInfo.Linkedin) {
		return ErrInvalidPayload
	}
	if !isValidURL(content.PortfolioInfo.Instagram) {
		return ErrInvalidPayload
	}
	for _, item := range content.Technical {
		if strings.TrimSpace(item.Title) == "" {
			return ErrInvalidPayload
		}
		if len(item.Description) > 2000 {
			return ErrInvalidPayload
		}
		if !isValidURL(item.Icon) {
			return ErrInvalidPayload
		}
	}
	for _, item := range content.Projects {
		if item.Index < 0 {
			return ErrInvalidPayload
		}
		if strings.TrimSpace(item.Title) == "" {
			return ErrInvalidPayload
		}
		if len(item.Description) > 3000 {
			return ErrInvalidPayload
		}
		if !isValidURL(item.RepoURL) {
			return ErrInvalidPayload
		}
		if !isValidURL(item.ProjectURL) {
			return ErrInvalidPayload
		}
		if !isValidURL(item.Image) {
			return ErrInvalidPayload
		}
		for _, imageURL := range item.Images {
			if !isValidURL(imageURL) {
				return ErrInvalidPayload
			}
		}
	}
	return nil
}

func validateTechnicalItem(item TechnicalItem) error {
	if item.Index < 0 {
		return ErrInvalidPayload
	}
	if strings.TrimSpace(item.Title) == "" {
		return ErrInvalidPayload
	}
	if len(item.Description) > 2000 {
		return ErrInvalidPayload
	}
	if !isValidURL(item.Icon) {
		return ErrInvalidPayload
	}
	return nil
}

func withContentIndexes(content ContentBody) ContentBody {
	if len(content.Technical) > 0 {
		normalizedTechnical := make([]TechnicalItem, len(content.Technical))
		copy(normalizedTechnical, content.Technical)
		for i := range normalizedTechnical {
			normalizedTechnical[i].Index = i
		}
		content.Technical = normalizedTechnical
	}

	if len(content.Projects) > 0 {
		normalizedProjects := make([]ProjectItem, len(content.Projects))
		copy(normalizedProjects, content.Projects)
		for i := range normalizedProjects {
			normalizedProjects[i].Index = i
		}
		content.Projects = normalizedProjects
	}

	return content
}
