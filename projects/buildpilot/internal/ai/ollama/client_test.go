package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ranjith-ka/buildpilot/internal/ai"
)

func TestStructuredDecodesModelContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["format"] == nil {
			t.Fatal("structured schema was not sent")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"role": "assistant", "content": `{"answer":"ok"}`},
		})
	}))
	defer server.Close()

	client := New(server.URL, "test-model", time.Second)
	var output struct {
		Answer string `json:"answer"`
	}
	err := client.Structured(context.Background(), ai.StructuredRequest{
		Prompt: "test", Schema: json.RawMessage(`{"type":"object"}`),
	}, &output)
	if err != nil {
		t.Fatalf("Structured() error = %v", err)
	}
	if output.Answer != "ok" {
		t.Fatalf("answer = %q", output.Answer)
	}
}

func TestStructuredFallsBackForLegacyOllamaFormat(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		var request struct {
			Format   json.RawMessage `json:"format"`
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if requests == 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"json: cannot unmarshal object into Go struct field ChatRequest.format of type string"}`))
			return
		}
		if string(request.Format) != `"json"` {
			t.Fatalf("fallback format = %s", request.Format)
		}
		if len(request.Messages) < 2 || !strings.Contains(request.Messages[1].Content, "Required JSON Schema") {
			t.Fatal("fallback prompt does not contain the schema")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"role": "assistant", "content": `{"answer":"fallback-ok"}`},
		})
	}))
	defer server.Close()

	client := New(server.URL, "test-model", time.Second)
	var output struct {
		Answer string `json:"answer"`
	}
	err := client.Structured(context.Background(), ai.StructuredRequest{
		Prompt: "test", Schema: json.RawMessage(`{"type":"object"}`),
	}, &output)
	if err != nil {
		t.Fatalf("Structured() error = %v", err)
	}
	if requests != 2 || output.Answer != "fallback-ok" {
		t.Fatalf("requests = %d, answer = %q", requests, output.Answer)
	}
}

func TestStructuredRejectsInvalidModelJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"role": "assistant", "content": "not-json"},
		})
	}))
	defer server.Close()

	client := New(server.URL, "test-model", time.Second)
	var output map[string]any
	if err := client.Structured(context.Background(), ai.StructuredRequest{}, &output); err == nil {
		t.Fatal("expected invalid structured response error")
	}
}
