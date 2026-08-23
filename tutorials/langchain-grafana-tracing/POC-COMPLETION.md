# LangGraph Grafana tracing POC: completion plan

## Current state

The repository already demonstrates a strong deterministic observability
workflow:

- typed LangGraph state, nodes, edges, and SQLite checkpoints;
- Ollama-based answer generation;
- local documentation retrieval;
- OpenTelemetry spans exported through a Collector to Tempo;
- trace-correlated structured logs in Loki;
- Tempo trace parsing and trace comparison;
- a Flask/UI demo, Docker Compose environment, and unit tests.

The current graph does **not** yet let the model choose investigation tools.
`graph.py` always runs:

```text
retrieve_documentation -> generate_answer -> persist_memory
```

The next POC milestone is therefore not more documentation retrieval. It is a
bounded, read-only tool-calling incident agent.

## Target POC

```text
Question + trace ID
       |
       v
Incident agent (Ollama)
       |
       +--> Tempo trace tool
       +--> Loki log tool
       +--> trace comparison tool
       |
       v
Evidence-backed RCA + confidence + next action
```

## Implementation order

### Milestone 1 — Tool wrappers

Create `incident_tools.py` with LangChain `@tool` wrappers around the existing
functions:

- `fetch_trace_from_tempo`;
- `fetch_logs_for_trace`;
- the existing trace comparison logic.

Return compact JSON-serializable evidence. Limit spans and log lines so a large
trace cannot flood the model context.

**Done when:** every tool has deterministic unit tests and useful descriptions
and rejects invalid trace IDs.

### Milestone 2 — Agent graph

Create `incident_graph.py` using:

- `MessagesState` or a custom extension;
- `ChatOllama.bind_tools(...)`;
- an agent/model node;
- `ToolNode`;
- `tools_condition`;
- the `tools -> agent` feedback edge;
- SQLite checkpointing by `thread_id`.

Keep the existing educational workflow intact so users can compare workflow and
agent designs.

**Done when:** the same compiled graph follows different tool paths for different
questions.

### Milestone 3 — Agent observability

Instrument:

- one root span per investigation;
- one span per model call;
- one span per tool call;
- tool name, duration, status, result size, model, and token usage;
- safe state-transition events.

Do not record raw prompts, log bodies, credentials, or full tool results by
default.

**Done when:** Grafana shows the complete reason/action/observation loop as a
waterfall, and Loki events share the same trace ID.

### Milestone 4 — Evaluation and guardrails

Add tests for:

- tool selection;
- no-tool final answers;
- multi-tool investigations;
- tool failure and timeout handling;
- maximum iteration/recursion limits;
- prompt injection inside log records;
- unsupported root-cause claims;
- sensitive-data redaction.

Use a deterministic fake tool-calling model in unit tests. Keep live Ollama tests
separate and optional.

**Done when:** CI can test graph behaviour without running Tempo, Loki, Grafana,
or Ollama.

### Milestone 5 — Demonstration workflow

Add a repeatable demo that creates or loads:

- one healthy trace;
- one slow database trace;
- one application-error trace;
- a baseline/candidate pair.

Show the selected tools, evidence, answer, confidence, and trace ID in the UI.

**Done when:** another developer can clone the repository and reproduce the demo
from the quickstart without manual data preparation.

## POC acceptance criteria

The POC is complete when all of these are true:

- The LLM selects from at least three read-only tools.
- At least two inputs produce different tool-call paths.
- Each conclusion includes observable evidence.
- The agent stops within a configured maximum number of model turns.
- Tool errors produce a controlled partial answer, not a crash or infinite retry.
- Unit tests do not require external services.
- A live Compose test reaches Tempo and Loki successfully.
- Grafana displays model and tool spans under one investigation trace.
- Secrets and raw sensitive payloads are absent from spans and logs.
- The README contains one command sequence for the complete demo.

## Recommended first vertical slice

Do only this before adding Kubernetes, GitHub, Prometheus, or deployment tools:

1. Wrap `fetch_trace_from_tempo` as `get_trace`.
2. Wrap `fetch_logs_for_trace` as `get_logs`.
3. Build the `agent -> tools -> agent` loop.
4. Add a maximum-turn limit.
5. Write one fake-model unit test proving the tool sequence.
6. Run one live Ollama investigation and inspect it in Grafana.

This is the smallest slice that proves the central POC claim: LangGraph can
dynamically investigate Grafana evidence, not merely generate text inside a
fixed workflow.
