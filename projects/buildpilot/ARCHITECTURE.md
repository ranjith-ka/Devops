# BuildPilot architecture

BuildPilot is a repository-to-Kubernetes build platform. A user connects a
repository, installs an outbound-only agent in a cluster, reviews an AI-created
build plan, and starts a rootless BuildKit build without sharing kubeconfig with
the control plane.

## System context

```text
GitHub App/webhooks                 Ollama (development)
        |                                  |
        v                                  v
  +---------------- BuildPilot control plane ----------------+
  | repository scanner | AI gateway | jobs | audit | web UI   |
  +-----------------------------------------------------------+
                            ^
                            | outbound HTTPS: poll, logs, status
                            v
                 +---- customer cluster agent ----+
                 | Kubernetes API | BuildKit jobs |
                 | cache PVC      | registry auth |
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

## Build state machine

```text
queued -> preparing -> building -> exporting -> deploying -> succeeded
                   \-> failed                \-> failed
       any active state -> cancelled
```

The control plane owns the desired state. The agent reports observed state with
monotonic sequence numbers so retries do not move a job backwards.

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
