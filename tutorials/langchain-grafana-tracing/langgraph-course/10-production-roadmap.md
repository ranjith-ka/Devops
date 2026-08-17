# Lesson 10: Production roadmap

## Learning objective

Know which tutorial choices must change as traffic, security requirements, and
team ownership grow.

## Current local choices

| Area | Tutorial implementation | Production direction |
| --- | --- | --- |
| API server | Flask development server | Production WSGI/ASGI workers |
| Checkpoints | SQLite named volume | Shared Postgres-backed checkpointer |
| Retrieval | In-memory lexical search | Evaluated hybrid/vector retrieval |
| Log delivery | Direct HTTP to Loki | stdout + Grafana Alloy/collector |
| Tempo storage | Local filesystem | Managed backend or object storage |
| Authentication | None | Identity, tenant isolation, authorization |
| Secrets | Environment variables | Secret manager and rotation |

## Scaling order

1. Add authentication and server-derived thread ownership.
2. Move checkpoints to shared durable storage.
3. Add bounded concurrency, timeouts, and retries.
4. Introduce production log shipping.
5. Evaluate retrieval quality before adopting a vector database.
6. Add sampling and retention policies for telemetry cost.
7. Load-test graph execution and failure recovery.

## Testing pyramid

- Unit tests: node state updates, parsers, retrieval ranking, diagnosis.
- Graph tests: routing, checkpoint restoration, retries, validation failures.
- Integration tests: Ollama, Tempo, Loki, and database adapters.
- End-to-end tests: browser question, follow-up memory, trace/log comparison.

## Final project

Extend the graph with classification and validation:

```text
START
  -> classify_request
      |- docs -> retrieve_documentation -+
      `- general ------------------------+
                                         v
                                  generate_answer
                                         |
                                  validate_answer
                                    |          |
                                  valid      retry/error
                                    |
                              persist_memory -> END
```

Acceptance criteria:

- every node has a stable span;
- every node has safe lifecycle logs;
- invalid output is never written to memory;
- routing is covered by deterministic tests;
- repeated attempts compare correctly in the UI;
- documentation and conversation memory remain separate.

Completing this project gives you the core architecture needed to add tools,
human approval, richer retrieval, and multi-step agent behavior safely.
