package ai

import (
	"context"
	"encoding/json"

	"github.com/ranjith-ka/buildpilot/internal/domain"
)

type StructuredRequest struct {
	SystemPrompt string
	Prompt       string
	Schema       json.RawMessage
}

type Provider interface {
	Structured(ctx context.Context, request StructuredRequest, output any) error
	Chat(ctx context.Context, messages []domain.Message) (string, error)
}
