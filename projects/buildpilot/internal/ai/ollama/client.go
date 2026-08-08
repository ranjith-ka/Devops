package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ranjith-ka/buildpilot/internal/ai"
	"github.com/ranjith-ka/buildpilot/internal/domain"
)

type Client struct {
	endpoint string
	model    string
	http     *http.Client
}

func New(endpoint, model string, timeout time.Duration) *Client {
	return &Client{
		endpoint: strings.TrimRight(endpoint, "/"),
		model:    model,
		http:     &http.Client{Timeout: timeout},
	}
}

type chatRequest struct {
	Model    string           `json:"model"`
	Messages []domain.Message `json:"messages"`
	Stream   bool             `json:"stream"`
	Format   json.RawMessage  `json:"format,omitempty"`
	Options  map[string]any   `json:"options,omitempty"`
}

type chatResponse struct {
	Message domain.Message `json:"message"`
	Error   string         `json:"error,omitempty"`
}

func (c *Client) Structured(ctx context.Context, request ai.StructuredRequest, output any) error {
	messages := []domain.Message{
		{Role: "system", Content: request.SystemPrompt},
		{Role: "user", Content: request.Prompt},
	}

	payload := chatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
		Format:   request.Schema,
		Options:  map[string]any{"temperature": 0},
	}
	content, err := c.complete(ctx, payload)
	if isLegacyFormatError(err) {
		// Older Ollama servers only accept "json" for format. Keep schema
		// validation instructions in the prompt and retry in JSON mode.
		payload.Format = json.RawMessage(`"json"`)
		payload.Messages[1].Content += "\n\nRequired JSON Schema:\n" + string(request.Schema)
		content, err = c.complete(ctx, payload)
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(content), output); err != nil {
		return fmt.Errorf("decode structured Ollama response: %w", err)
	}
	return nil
}

func isLegacyFormatError(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	return strings.Contains(message, "cannot unmarshal object") &&
		strings.Contains(message, "ChatRequest.format")
}

func (c *Client) Chat(ctx context.Context, messages []domain.Message) (string, error) {
	return c.complete(ctx, chatRequest{Model: c.model, Messages: messages, Stream: false})
}

func (c *Client) complete(ctx context.Context, payload chatRequest) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("encode Ollama request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create Ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("call Ollama: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", fmt.Errorf("read Ollama response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Ollama returned %s: %s", resp.Status, strings.TrimSpace(string(data)))
	}
	var result chatResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("decode Ollama response: %w", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("Ollama: %s", result.Error)
	}
	return result.Message.Content, nil
}
