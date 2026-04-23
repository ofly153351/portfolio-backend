package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Repository interface {
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
	StreamChat(ctx context.Context, req ChatRequest, onEvent func(StreamEvent) error) error
}

type AIServiceRepository struct {
	baseURL    string
	httpClient *http.Client
}

func NewAIServiceRepository(baseURL string) *AIServiceRepository {
	return NewAIServiceRepositoryWithTimeout(baseURL, 45*time.Second)
}

func NewAIServiceRepositoryWithTimeout(baseURL string, timeout time.Duration) *AIServiceRepository {
	if timeout <= 0 {
		timeout = 45 * time.Second
	}
	return &AIServiceRepository{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *AIServiceRepository) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	if r.baseURL == "" {
		return ChatResponse{}, ErrAIServiceURLMissing
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return ChatResponse{}, err
	}

	url := r.baseURL + "/chat"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return ChatResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		return ChatResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ChatResponse{}, fmt.Errorf("ai service error (%d): %s", resp.StatusCode, string(respBody))
	}

	var out ChatResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return ChatResponse{}, err
	}
	return out, nil
}

func (r *AIServiceRepository) StreamChat(ctx context.Context, req ChatRequest, onEvent func(StreamEvent) error) error {
	if r.baseURL == "" {
		return ErrAIServiceURLMissing
	}
	if onEvent == nil {
		return fmt.Errorf("onEvent callback is required")
	}

	resp, err := r.Chat(ctx, req)
	if err != nil {
		return err
	}

	for _, token := range strings.Fields(resp.Answer) {
		if err := onEvent(StreamEvent{Type: "token", Token: token + " "}); err != nil {
			return err
		}
	}
	if err := onEvent(StreamEvent{
		Type:      "done",
		SessionID: resp.SessionID,
		Provider:  resp.Provider,
		Usage:     &resp.Usage,
	}); err != nil {
		return err
	}
	return nil
}
