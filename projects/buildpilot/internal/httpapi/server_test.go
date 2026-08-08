package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ranjith-ka/buildpilot/internal/ai"
	"github.com/ranjith-ka/buildpilot/internal/domain"
	"github.com/ranjith-ka/buildpilot/internal/store"
)

type fakeProvider struct{}

func (fakeProvider) Structured(_ context.Context, _ ai.StructuredRequest, output any) error {
	encoded, _ := json.Marshal(domain.BuildPlan{Summary: "test", Builder: "buildkit"})
	return json.Unmarshal(encoded, output)
}

func TestCreateBuildQueuesValidatedDockerfile(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Dockerfile"), []byte("FROM scratch\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	memory := store.NewMemory()
	server, err := New(slog.New(slog.NewTextHandler(io.Discard, nil)), fakeProvider{}, memory, root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/builds", strings.NewReader(`{
		"repository":".", "dockerfile":"Dockerfile", "image":"example/api:dev", "agent_id":"local-kind"
	}`))
	request.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	job, found := memory.NextJob("local-kind")
	if !found || job.Image != "example/api:dev" || job.Metadata["dockerfile"] != "Dockerfile" {
		t.Fatalf("queued job = %#v, found = %v", job, found)
	}
}

func TestCreateBuildRejectsMissingDockerfile(t *testing.T) {
	root := t.TempDir()
	server, err := New(slog.New(slog.NewTextHandler(io.Discard, nil)), fakeProvider{}, store.NewMemory(), root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/builds", strings.NewReader(`{
		"repository":".", "dockerfile":"Dockerfile", "image":"example/api:dev", "agent_id":"local-kind"
	}`))
	request.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func (fakeProvider) Chat(context.Context, []domain.Message) (string, error) { return "", nil }

func TestHealth(t *testing.T) {
	server, err := New(slog.New(slog.NewTextHandler(io.Discard, nil)), fakeProvider{}, store.NewMemory(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
}

func TestAnalyzeRejectsPathOutsideRoot(t *testing.T) {
	root := t.TempDir()
	server, err := New(slog.New(slog.NewTextHandler(io.Discard, nil)), fakeProvider{}, store.NewMemory(), root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/repositories/analyze", strings.NewReader(`{"path":"../outside"}`))
	request.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
