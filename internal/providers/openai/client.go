package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/toshikimiyagawa/codex-approval-gate/internal/judge"
)

type Config struct {
	BaseURL string
	Model   string
	APIKey  string
	Timeout time.Duration
}

type Client struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

func New(cfg Config) Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		model:   cfg.Model,
		apiKey:  cfg.APIKey,
		client:  &http.Client{Timeout: timeout},
	}
}

func (c Client) Decide(ctx context.Context, req judge.ProviderRequest) (judge.ProviderResponse, error) {
	body := chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: req.SystemPrompt},
			{Role: "user", Content: req.UserPrompt},
		},
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return judge.ProviderResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(encoded))
	if err != nil {
		return judge.ProviderResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return judge.ProviderResponse{}, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return judge.ProviderResponse{}, fmt.Errorf("openai-compatible provider returned status %d", httpResp.StatusCode)
	}

	var completion chatCompletionResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&completion); err != nil {
		return judge.ProviderResponse{}, err
	}
	if len(completion.Choices) == 0 {
		return judge.ProviderResponse{}, fmt.Errorf("openai-compatible provider returned no choices")
	}
	return judge.ParseProviderResponse(completion.Choices[0].Message.Content)
}

type chatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}
