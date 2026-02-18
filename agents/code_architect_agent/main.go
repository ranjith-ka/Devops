package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const baseSystemPrompt = `You are Code Architect Partner, an expert AI software engineer helping with:
1) code editing
2) architecture design and refactoring
3) implementation planning with trade-offs

Rules:
- Keep responses concise and action-oriented.
- For code changes, return unified diff snippets first, then a short rationale.
- Prefer minimal, safe, incremental changes.
- If information is missing, ask one focused question.
- When discussing architecture, include: context, options, recommendation, and migration steps.
- Never invent existing files/functions; mark proposals clearly.`

type chatRequest struct {
	Message string `json:"message"`
	Mode    string `json:"mode,omitempty"`
	Context string `json:"context,omitempty"`
}

type chatResponse struct {
	Reply string `json:"reply"`
}

type llmClient struct {
	httpClient *http.Client
	github     bool

	githubToken    string
	githubModel    string
	githubEndpoint string
}

func newLLMClient() (*llmClient, error) {
	if token := os.Getenv("GITHUB_MODELS_TOKEN"); token != "" {
		model := os.Getenv("GITHUB_MODELS_MODEL")
		if model == "" {
			model = "openai/gpt-5-mini"
		}

		endpoint := os.Getenv("GITHUB_MODELS_ENDPOINT")
		if endpoint == "" {
			endpoint = "https://models.github.ai/inference/chat/completions"
		}

		return &llmClient{
			httpClient:     &http.Client{Timeout: 90 * time.Second},
			github:         true,
			githubToken:    token,
			githubModel:    model,
			githubEndpoint: strings.TrimRight(endpoint, "/"),
		}, nil
	}

	return nil, errors.New("missing model configuration: set GITHUB_MODELS_TOKEN (and optional GITHUB_MODELS_MODEL/GITHUB_MODELS_ENDPOINT)")
}

func (c *llmClient) chat(ctx context.Context, mode, userContext, message string) (string, error) {
	systemPrompt := buildPrompt(mode, userContext)

	requestBody := map[string]any{
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": message,
			},
		},
	}

	endpoint := c.githubEndpoint
	requestBody["model"] = c.githubModel

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.githubToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("llm request failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) == 0 {
		return "", errors.New("llm returned empty choices")
	}

	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}

func buildPrompt(mode, contextText string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	var modePrompt string

	switch mode {
	case "edit":
		modePrompt = "Focus on code edits. Output should begin with minimal unified diff snippets."
	case "architecture", "arch":
		modePrompt = "Focus on architecture quality: trade-offs, target design, migration steps, and risks."
	default:
		modePrompt = "Handle both implementation details and architecture concerns, keeping suggestions incremental."
	}

	if strings.TrimSpace(contextText) == "" {
		return baseSystemPrompt + "\n\nMode guidance: " + modePrompt
	}

	return baseSystemPrompt + "\n\nMode guidance: " + modePrompt + "\n\nWorkspace context:\n" + contextText
}

func main() {
	loadEnvFiles()

	client, err := newLLMClient()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request chatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(request.Message) == "" {
			http.Error(w, "message is required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 100*time.Second)
		defer cancel()

		reply, err := client.chat(ctx, request.Mode, request.Context, request.Message)
		if err != nil {
			log.Printf("chat error: %v", err)
			http.Error(w, "failed to generate response", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(chatResponse{Reply: reply}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	addr := ":" + port
	log.Printf("Code Architect Agent running on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func loadEnvFiles() {
	files := []string{".env", filepath.Join("agents", "code_architect_agent", ".env")}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			idx := strings.Index(line, "=")
			if idx <= 0 {
				continue
			}

			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			value = strings.Trim(value, `"`)
			value = strings.Trim(value, `'`)

			if key == "" {
				continue
			}

			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, value)
			}
		}
	}
}
