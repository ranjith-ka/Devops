# BuildPilot architecture

BuildPilot is an AI-oriented remote verification platform. A coding agent or
user submits an immutable source revision and optional patch, and an
outbound-only cluster agent runs compilation, tests, analysis, container
building, and integration tests without sharing kubeconfig with the control
plane. Rootless BuildKit is the image-building worker, not the product boundary.

## System context

```text
GitHub App/webhooks                 Ollama (development)
        |                                  |
        v                                  v
  +---------------- BuildPilot control plane ----------------+
  | verification API | scheduler | results | audit | web UI  |
  +-----------------------------------------------------------+
                            ^
                            | outbound HTTPS: poll, logs, status
                            v
                 +---- customer cluster agent ----+
                 | .NET Jobs | Sonar | BuildKit   |
                 | caches | reports | registry    |
                 +-------------------------------+
```

## Trust boundaries

- Repository content, build logs, Dockerfiles, and model output are untrusted.
- The control plane never accepts or stores customer kubeconfig files.
- The agent owns cluster credentials and initiates all control-plane traffic.
- Registry credentials remain in the customer cluster.
- AI output is advisory. File writes, builds, and deployments require a policy
  decision and, in the MVP, explicit user approval.

## Go module layout

```text
cmd/control-plane       HTTP API and embedded web UI
cmd/cluster-agent       outbound agent process
internal/ai             provider-neutral model API
internal/ai/ollama      Ollama implementation
internal/analyzer       deterministic repository inspection
internal/domain         API and job domain types
internal/httpapi        routes, validation, and response encoding
internal/store          persistence interfaces and memory implementation
web                     landing page assets
```

## Verification state machine

```text
queued -> provisioning -> restoring -> compiling -> testing
                                                 |         |
                                                 | fast    | full
                                                 v         v
                                           succeeded   analyzing
                                                          |
                                                     image_building
                                                          |
                                                   integration_testing
                                                          |
                                                     succeeded

any active state -> failed | cancelled | superseded | timed_out
```

The control plane and its durable database own desired and recorded state.
Kubernetes resources are execution state, not the system of record. The agent
reports observed state with monotonic sequence numbers so retries do not move a
verification backwards.

## AI boundary

The deterministic scanner parses known files first. Only a bounded summary and
selected sanitized files are provided to AI. Ollama must return JSON matching a
schema for repository plans and failure diagnoses. Provider output is validated
before it enters the domain model.

Initial AI use cases:

1. Explain detected project structure.
2. Recommend a builder, target platforms, port, and health endpoint.
3. Diagnose a bounded build-log excerpt.
4. Propose file changes for user review.

## Production evolution

- Replace the memory store with PostgreSQL and an outbox-backed job queue.
- Add GitHub App installation tokens and verified webhook ingestion.
- Authenticate agents using one-time enrollment followed by rotated mTLS.
- Replace polling with a reconnecting bidirectional stream when scale requires.
- Use object storage for source bundles and complete log archives.
- Add a BuildKit controller using `client-go`; keep build credentials mounted
  from cluster-side Secrets.
- Add per-tenant quotas, admission policies, audit events, and billing meters.
