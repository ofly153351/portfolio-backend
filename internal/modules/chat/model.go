package chat

type ChatRequest struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id"`
	TopK      int    `json:"top_k"`
	Lang      string `json:"lang"`
}

type ChatSource struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Score float64 `json:"score"`
}

type ChatResponse struct {
	Answer    string       `json:"answer"`
	Sources   []ChatSource `json:"sources"`
	SessionID string       `json:"session_id"`
	Provider  string       `json:"provider,omitempty"`
	Usage     ChatUsage    `json:"usage"`
}

type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamEvent struct {
	Type      string     `json:"type"`
	Token     string     `json:"token,omitempty"`
	Message   string     `json:"message,omitempty"`
	Error     string     `json:"error,omitempty"`
	SessionID string     `json:"session_id,omitempty"`
	Provider  string     `json:"provider,omitempty"`
	Usage     *ChatUsage `json:"usage,omitempty"`
}
