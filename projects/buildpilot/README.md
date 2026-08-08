# BuildPilot

BuildPilot is an early product scaffold for AI-assisted remote container builds.
It combines deterministic repository inspection, local Ollama analysis, and an
outbound cluster-agent architecture. The current phase is deliberately safe:
AI returns reviewable plans and diagnoses but cannot write files or execute
cluster operations.

Read [ARCHITECTURE.md](ARCHITECTURE.md) for the system design and
[PLAN.md](PLAN.md) for delivery phases.

## Requirements

- Go 1.24 or newer
- Ollama
- A model with structured-output support

## Start Ollama

```bash
ollama serve
ollama pull qwen2.5-coder:7b
```

## Start the control plane

From this directory:

```bash
REPOSITORY_ROOT=../.. \
OLLAMA_MODEL=qwen2.5-coder:7b \
go run ./cmd/control-plane
```

Or start it with the configured local defaults:

```bash
make start
```

You can override any setting when needed:

```bash
make start REPOSITORY_ROOT=/path/to/repos OLLAMA_MODEL=qwen2.5-coder:7b
```

`make run` remains an alias for `make start`. Run the test suite with
`make test`.

Open <http://localhost:8090>. The primary UI flow validates a repository,
Dockerfile, image name, and agent ID, then queues the build without calling AI.
The repository path is resolved under `REPOSITORY_ROOT`; paths outside that
root are rejected. Ollama is used only by the optional failure-diagnosis flow.

The phase-0 cluster agent can receive this queued job but does not execute it
yet. Source transfer and the rootless BuildKit Kubernetes Job executor are the
next implementation milestone; a queued response must not be treated as a
successful image build.

## API examples

Analyze the .NET Skaffold example:

```bash
curl http://localhost:8090/api/v1/repositories/analyze \
  -H 'Content-Type: application/json' \
  -d '{"path":"examples/dotnet-skaffold"}'
```

Diagnose a build error:

```bash
curl http://localhost:8090/api/v1/builds/diagnose \
  -H 'Content-Type: application/json' \
  -d '{"logs":"Requested SDK 10.0.107; installed SDK 10.0.302"}'
```

## Start the phase-0 agent

The current agent proves the outbound heartbeat and job-polling boundary; it
does not yet create BuildKit Jobs.

```bash
AGENT_ID=local-kind \
CLUSTER_NAME=kind-k8s \
CONTROL_PLANE_URL=http://localhost:8090 \
go run ./cmd/cluster-agent
```

## Safety notes

- Use a local test repository root; do not point this prototype at sensitive
  directories.
- Repository text and logs are untrusted and may contain prompt injection.
- Basic redaction is illustrative, not sufficient for production.
- Authentication, tenant isolation, durable storage, signed agent requests,
  GitHub App integration, and Kubernetes execution belong to subsequent phases.
