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

const tutorialSystemPrompt = `You are Tutorial Mentor Agent, an expert technical tutor.

Your responsibilities:
1) create structured tutorials from beginner to advanced
2) explain concepts with practical examples
3) generate hands-on exercises and checkpoints
4) adapt difficulty based on user input

Response style rules:
- Keep explanations clear and concise.
- Use step-by-step sections when teaching.
- Include prerequisites, learning goals, and expected outcomes.
- For coding tutorials, provide runnable examples and verification steps.
- If user goal is vague, ask one focused clarification question.
- Do not output unsafe or destructive commands.`

type tutorialRequest struct {
	Message    string `json:"message"`
	Mode       string `json:"mode,omitempty"`
	Topic      string `json:"topic,omitempty"`
	Level      string `json:"level,omitempty"`
	OutputFile string `json:"outputFile,omitempty"`
}

type tutorialResponse struct {
	Reply     string `json:"reply"`
	FilePath  string `json:"filePath"`
	FileName  string `json:"fileName"`
	SavedIn   string `json:"savedIn"`
	Timestamp string `json:"timestamp"`
}

type llmClient struct {
	httpClient     *http.Client
	githubToken    string
	githubModel    string
	githubEndpoint string
}

func newLLMClient() (*llmClient, error) {
	token := os.Getenv("GITHUB_MODELS_TOKEN")
	if token == "" {
		return nil, errors.New("missing GITHUB_MODELS_TOKEN")
	}

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
		githubToken:    token,
		githubModel:    model,
		githubEndpoint: strings.TrimRight(endpoint, "/"),
	}, nil
}

func (c *llmClient) chat(ctx context.Context, mode, topic, level, message string) (string, error) {
	systemPrompt := buildTutorialPrompt(mode, topic, level)

	requestBody := map[string]any{
		"model": c.githubModel,
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

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.githubEndpoint, bytes.NewReader(bodyBytes))
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

func buildTutorialPrompt(mode, topic, level string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	level = strings.ToLower(strings.TrimSpace(level))
	if level == "" {
		level = "intermediate"
	}

	var modeGuide string
	switch mode {
	case "plan":
		modeGuide = "Generate a learning path with milestones and weekly goals."
	case "lesson":
		modeGuide = "Generate a focused lesson with examples and practice tasks."
	case "exercise":
		modeGuide = "Generate practical exercises with expected output and hints."
	default:
		modeGuide = "Provide a balanced tutorial response with explanation, steps, and practice."
	}

	var builder strings.Builder
	builder.WriteString(tutorialSystemPrompt)
	builder.WriteString("\n\nMode guidance: ")
	builder.WriteString(modeGuide)
	builder.WriteString("\nTarget level: ")
	builder.WriteString(level)

	if strings.TrimSpace(topic) != "" {
		builder.WriteString("\nTopic: ")
		builder.WriteString(strings.TrimSpace(topic))
	}

	return builder.String()
}

func loadEnvFiles() {
	files := []string{".env", filepath.Join("agents", "tutorial_agent", ".env")}

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

func main() {
	loadEnvFiles()

	client, err := newLLMClient()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/tutorial", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request tutorialRequest
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

		reply, err := client.chat(ctx, request.Mode, request.Topic, request.Level, request.Message)
		if err != nil {
			log.Printf("tutorial error: %v", err)
			http.Error(w, "failed to generate response", http.StatusBadGateway)
			return
		}

		fullPath, fileName, outDir, err := saveTutorialMarkdown(request, reply)
		if err != nil {
			log.Printf("save markdown error: %v", err)
			http.Error(w, "failed to save markdown", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tutorialResponse{
			Reply:     reply,
			FilePath:  fullPath,
			FileName:  fileName,
			SavedIn:   outDir,
			Timestamp: time.Now().Format(time.RFC3339),
		}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	addr := ":" + port
	log.Printf("Tutorial Agent running on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func saveTutorialMarkdown(req tutorialRequest, reply string) (string, string, string, error) {
	outDir, err := resolveOutputDir()
	if err != nil {
		return "", "", "", err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", "", "", err
	}

	fileName := buildOutputFileName(req)
	fullPath := filepath.Join(outDir, fileName)

	content := buildMarkdownContent(req, reply)
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return "", "", "", err
	}

	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		absPath = fullPath
	}

	return absPath, fileName, outDir, nil
}

func resolveOutputDir() (string, error) {
	if configured := strings.TrimSpace(os.Getenv("TUTORIAL_OUTPUT_DIR")); configured != "" {
		return configured, nil
	}

	baseCandidates := []string{"minikube", filepath.Join("..", "..", "minikube")}
	for _, base := range baseCandidates {
		info, err := os.Stat(base)
		if err == nil && info.IsDir() {
			return filepath.Join(base, "tutorial-ai"), nil
		}
	}

	return "", errors.New("could not locate minikube folder; set TUTORIAL_OUTPUT_DIR")
}

func buildOutputFileName(req tutorialRequest) string {
	if custom := strings.TrimSpace(req.OutputFile); custom != "" {
		name := filepath.Base(custom)
		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			name += ".md"
		}
		return name
	}

	topic := sanitizeFilePart(req.Topic)
	if topic == "" {
		topic = "tutorial"
	}

	mode := sanitizeFilePart(req.Mode)
	if mode == "" {
		mode = "general"
	}

	stamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-%s-%s.md", topic, mode, stamp)
}

func sanitizeFilePart(input string) string {
	clean := strings.TrimSpace(strings.ToLower(input))
	if clean == "" {
		return ""
	}

	var builder strings.Builder
	lastDash := false
	for _, r := range clean {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastDash = false
			continue
		}

		if !lastDash {
			builder.WriteRune('-')
			lastDash = true
		}
	}

	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return ""
	}

	return result
}

func buildMarkdownContent(req tutorialRequest, reply string) string {
	now := time.Now().Format(time.RFC3339)

	var b strings.Builder
	b.WriteString("# Tutorial Response\n\n")
	b.WriteString("- Generated: ")
	b.WriteString(now)
	b.WriteString("\n")

	if strings.TrimSpace(req.Topic) != "" {
		b.WriteString("- Topic: ")
		b.WriteString(strings.TrimSpace(req.Topic))
		b.WriteString("\n")
	}

	if strings.TrimSpace(req.Level) != "" {
		b.WriteString("- Level: ")
		b.WriteString(strings.TrimSpace(req.Level))
		b.WriteString("\n")
	}

	if strings.TrimSpace(req.Mode) != "" {
		b.WriteString("- Mode: ")
		b.WriteString(strings.TrimSpace(req.Mode))
		b.WriteString("\n")
	}

	b.WriteString("\n## Request\n\n")
	b.WriteString(req.Message)
	b.WriteString("\n\n## Tutorial\n\n")
	b.WriteString(reply)
	b.WriteString("\n")

	return b.String()
}
