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
- Accept an immutable base commit plus an optional patch.
- Run `dotnet restore`, compile, and unit tests in an isolated Kubernetes Job.
- Parse compiler and TRX output into structured, agent-readable diagnostics.
- Stream logs and support timeout, cancellation, and superseding older work.
- Upload test reports and logs before automatic resource cleanup.

Exit criterion: an agent can submit a failing .NET patch, receive an actionable
file-and-line diagnosis, correct it, and pass fast verification in one connected
cluster without sharing kubeconfig.

## Phase 2: full verification and container output

- Add full test, coverage, and Sonar quality-gate stages.
- Create rootless BuildKit Jobs with persistent and registry caches.
- Return the immutable image digest, SBOM, provenance, and report links.
- Deploy by digest into a temporary namespace and run integration tests.
- Distinguish retryable infrastructure failures from code and policy failures.
- Package production installation as a versioned Helm chart and signed images.

Exit criterion: a verified .NET change produces a tested immutable image digest
and complete evidence, and failed quality gates prevent promotion.

## Phase 3: GitHub onboarding

- Register a GitHub App with minimal repository permissions.
- Add installation callback and repository selection.
- Verify and persist push and pull-request webhooks.
- Publish GitHub Checks with logs and AI diagnosis.
- Generate patches as reviewable suggestions; never commit without approval.

Exit criterion: onboarding from landing page to first build takes under ten
minutes and requires no hand-written pipeline file.

## Phase 4: deployment and teams

- Helm and raw Kubernetes deployment plans.
- Environment promotion, rollout status, and rollback.
- Organizations, roles, audit logs, quotas, and retention.
- OIDC/workload identity for ECR, ACR, and Artifact Registry.
- AMD64, ARM64, and multi-platform build matrices.

## Non-goals for the MVP

- General-purpose workflow automation.
- Arbitrary AI tool execution.
- Production cluster administration.
- Replacing GitHub Actions as a general-purpose CI system.
- Building a new workflow engine, registry, object store, or Sonar replacement.
