package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ranjith-ka/Devops/services/inference/internal/ai"
)

type handlers struct {
	modelRouter *ai.Router
}

func newHandlers() *handlers { return &handlers{modelRouter: ai.NewRouter()} }

type response struct {
	RequestID       string    `json:"request_id"`
	Model           string    `json:"model,omitempty"`
	PromptTemplate  string    `json:"prompt_template,omitempty"`
	Summary         string    `json:"summary"`
	Confidence      float64   `json:"confidence,omitempty"`
	Evidence        []string  `json:"evidence,omitempty"`
	Recommendations []string  `json:"recommendations,omitempty"`
	Artifacts       []string  `json:"artifacts,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
}

type requestEnvelope struct {
	OrgID      string         `json:"org_id"`
	Repository string         `json:"repository"`
	Pipeline   string         `json:"pipeline"`
	Deployment string         `json:"deployment"`
	Prompt     string         `json:"prompt"`
	Context    requestContext `json:"context"`
}

type requestContext struct {
	Branch      string `json:"branch"`
	Environment string `json:"environment"`
	TimeRange   string `json:"time_range"`
}

func (h *handlers) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *handlers) chat(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "chat", "Use the evidence drawer to answer with citations and next actions.")
}
func (h *handlers) analyzePipeline(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "analyze-pipeline", "Pipeline analysis ready with bottlenecks, flaky tests, and optimization opportunities.")
}
func (h *handlers) deploymentSummary(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "deployment-summary", "Deployment summary assembled from run history, logs, and telemetry.")
}
func (h *handlers) rootCause(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "root-cause", "Root cause hypotheses ranked by confidence and evidence density.")
}
func (h *handlers) optimize(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "optimize", "Optimization plan generated for cost, speed, and reliability improvements.")
}
func (h *handlers) securityReview(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "security-review", "Security review completed with policy gaps and remediation guidance.")
}
func (h *handlers) explainWorkflow(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "explain-workflow", "Workflow explanation generated with stages, gates, and failure points.")
}
func (h *handlers) generatePipeline(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "generate-pipeline", "Generated pipeline skeleton with approval gates and deployment stages.")
}
func (h *handlers) generateTerraform(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "generate-terraform", "Generated Terraform scaffold for reproducible infrastructure.")
}
func (h *handlers) generateKubernetes(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "generate-kubernetes", "Generated Kubernetes manifests with secure defaults and rollout hooks.")
}
func (h *handlers) generateVideoPlan(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "generate-video-plan", "Generated YouTube content brief with title, hook, segments, CTA, and asset checklist.")
}
func (h *handlers) uploadVideo(w http.ResponseWriter, r *http.Request) {
	h.writeAIResponse(w, r, "upload-video", "Prepared YouTube upload workflow with metadata, timing, and publishing safeguards.")
}

func (h *handlers) writeAIResponse(w http.ResponseWriter, r *http.Request, task string, summary string) {
	var req requestEnvelope
	_ = json.NewDecoder(r.Body).Decode(&req)
	providerPlan := h.modelRouter.Plan(task, req.Prompt)
	writeJSON(w, http.StatusOK, response{
		RequestID:       "request-" + task,
		Model:           providerPlan.Model,
		PromptTemplate:  providerPlan.PromptTemplate,
		Summary:         summary,
		Confidence:      providerPlan.Confidence,
		Evidence:        []string{"logs", "metrics", "deployment history", "PR diff"},
		Recommendations: []string{"inspect the latest failing stage", "review the generated evidence bundle", "apply the suggested remediation"},
		Artifacts:       []string{"summary.md", "yaml-diff", "risk-score.json"},
		Timestamp:       time.Now().UTC(),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
