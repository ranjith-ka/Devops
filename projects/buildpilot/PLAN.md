# BuildPilot product plan

## Phase 0: executable prototype

- [x] Independent Go module.
- [x] Landing page and health endpoint.
- [x] Deterministic local repository scanner.
- [x] Provider-neutral AI interface and Ollama client.
- [x] Structured repository analysis and build-log diagnosis endpoints.
- [x] Outbound cluster-agent heartbeat and polling skeleton.
- [ ] Connect the agent to Kubernetes and create BuildKit Jobs.

Exit criterion: a developer can run Ollama and the control plane, analyze a
local repository from the browser, and receive a structured recommendation.

## Phase 1: single-user build loop

- Persist repositories, agents, build plans, and jobs in PostgreSQL.
- Add one-time agent enrollment and signed requests.
- Package the agent as a Helm chart with namespace-scoped RBAC.
- Upload a sanitized source bundle to object storage.
- Create rootless BuildKit Jobs with persistent and registry caches.
- Stream logs and support cancellation.
- Support build-only, OCI download, and registry-push outputs.

Exit criterion: one private repository can be built reliably in one connected
cluster without sharing kubeconfig.

## Phase 2: GitHub onboarding

- Register a GitHub App with minimal repository permissions.
- Add installation callback and repository selection.
- Verify and persist push and pull-request webhooks.
- Publish GitHub Checks with logs and AI diagnosis.
- Generate patches as reviewable suggestions; never commit without approval.

Exit criterion: onboarding from landing page to first build takes under ten
minutes and requires no hand-written pipeline file.

## Phase 3: deployment and teams

- Helm and raw Kubernetes deployment plans.
- Environment promotion, rollout status, and rollback.
- Organizations, roles, audit logs, quotas, and retention.
- OIDC/workload identity for ECR, ACR, and Artifact Registry.
- AMD64, ARM64, and multi-platform build matrices.

## Non-goals for the MVP

- General-purpose workflow automation.
- Arbitrary AI tool execution.
- Production cluster administration.
- Replacing GitHub Actions for tests unrelated to container delivery.
