package ai

import "testing"

func TestPlanSelectsSecurityTemplate(t *testing.T) {
	plan := NewRouter().Plan("security-review", "scan the release for vulnerabilities")

	if plan.Provider != "claude" {
		t.Fatalf("expected claude provider, got %s", plan.Provider)
	}
	if plan.PromptTemplate != "security-review-v1" {
		t.Fatalf("expected security-review-v1 template, got %s", plan.PromptTemplate)
	}
}

func TestPlanSelectsArtifactTemplate(t *testing.T) {
	plan := NewRouter().Plan("generate-kubernetes", "build a deployment yaml")

	if plan.Provider != "gpt" {
		t.Fatalf("expected gpt provider, got %s", plan.Provider)
	}
	if plan.PromptTemplate != "delivery-artifact-v1" {
		t.Fatalf("expected delivery-artifact-v1 template, got %s", plan.PromptTemplate)
	}
}

func TestPlanSelectsYouTubeAutomationTemplate(t *testing.T) {
	plan := NewRouter().Plan("upload-video", "create a youtube shorts upload package with title and thumbnail")

	if plan.Provider != "gpt" {
		t.Fatalf("expected gpt provider, got %s", plan.Provider)
	}
	if plan.PromptTemplate != "youtube-automation-v1" {
		t.Fatalf("expected youtube-automation-v1 template, got %s", plan.PromptTemplate)
	}
}
