package domain

import "time"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RepositoryFacts struct {
	Root            string   `json:"root"`
	Files           []string `json:"files"`
	Languages       []string `json:"languages"`
	ProjectFiles    []string `json:"project_files"`
	Dockerfiles     []string `json:"dockerfiles"`
	KubernetesFiles []string `json:"kubernetes_files"`
	HasSkaffold     bool     `json:"has_skaffold"`
}

type BuildPlan struct {
	Summary           string   `json:"summary"`
	Language          string   `json:"language"`
	Framework         string   `json:"framework"`
	Builder           string   `json:"builder"`
	Platforms         []string `json:"platforms"`
	ContainerPort     int      `json:"container_port"`
	HealthEndpoint    string   `json:"health_endpoint"`
	RequiredFiles     []string `json:"required_files"`
	Warnings          []string `json:"warnings"`
	NeedsUserApproval bool     `json:"needs_user_approval"`
}

type BuildDiagnosis struct {
	Summary          string   `json:"summary"`
	Category         string   `json:"category"`
	Confidence       float64  `json:"confidence"`
	Evidence         []string `json:"evidence"`
	SuggestedChanges []string `json:"suggested_changes"`
	SafeToAutoApply  bool     `json:"safe_to_auto_apply"`
}

type Agent struct {
	ID            string    `json:"id"`
	ClusterName   string    `json:"cluster_name"`
	Architecture  string    `json:"architecture"`
	Version       string    `json:"version"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

type JobStatus string

const (
	JobQueued    JobStatus = "queued"
	JobPreparing JobStatus = "preparing"
	JobBuilding  JobStatus = "building"
	JobSucceeded JobStatus = "succeeded"
	JobFailed    JobStatus = "failed"
)

type Job struct {
	ID        string            `json:"id"`
	AgentID   string            `json:"agent_id"`
	Image     string            `json:"image"`
	CommitSHA string            `json:"commit_sha"`
	Status    JobStatus         `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}
