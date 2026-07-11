package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatIncludesPromptTemplate(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(`{"org_id":"org-1","prompt":"summarize the deployment"}`))
	rr := httptest.NewRecorder()

	newHandlers().chat(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp["prompt_template"] == "" {
		t.Fatalf("expected prompt_template in response")
	}
	if resp["summary"] == "" {
		t.Fatalf("expected summary in response")
	}
}

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	newHandlers().healthz(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGenerateVideoPlanIncludesPromptTemplate(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/generate-video-plan", strings.NewReader(`{"org_id":"org-1","prompt":"create a youtube shorts outline and upload plan"}`))
	rr := httptest.NewRecorder()

	newHandlers().generateVideoPlan(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp["prompt_template"] != "youtube-automation-v1" {
		t.Fatalf("expected youtube-automation-v1 template, got %v", resp["prompt_template"])
	}
}
