# LangGraph observability application: step-by-step course

This course teaches the architecture implemented in this repository. It starts
with graph state and finishes with adding production-ready nodes that appear in
Tempo waterfalls and Loki logs automatically.

## How to study

1. Keep the application running with `docker compose up -d --build`.
2. Read lessons in order.
3. Open the linked source file while reading.
4. Complete each exercise before continuing.
5. Run `.venv/bin/python -m unittest discover -s tests -v` after code changes.

## Curriculum

| Lesson | Topic | Main code |
| --- | --- | --- |
| [01](./01-architecture-tour.md) | System architecture and request lifecycle | [`ui.py`](../ui.py), [`graph.py`](../graph.py) |
| [02](./02-state-nodes-edges.md) | Typed state, nodes, and edges | [`graph.py`](../graph.py) |
| [03](./03-invocation-and-thread-identity.md) | API invocation and thread identity | [`ui.py`](../ui.py) |
| [04](./04-checkpoints-and-memory.md) | SQLite checkpoints and durable memory | [`graph.py`](../graph.py) |
| [05](./05-documentation-retrieval.md) | Grounding answers in local documentation | [`documentation.py`](../documentation.py) |
| [06](./06-tempo-node-tracing.md) | OpenTelemetry and Tempo node spans | [`tracing.py`](../tracing.py), [`graph.py`](../graph.py) |
| [07](./07-loki-log-correlation.md) | Structured Loki logs and trace correlation | [`loki.py`](../loki.py) |
| [08](./08-trace-and-log-comparison.md) | Real trace comparison and waterfall UI | [`tempo.py`](../tempo.py), [`trace_analyzer.py`](../trace_analyzer.py) |
| [09](./09-extending-the-graph.md) | Add tools, routing, validation, and retries | [`graph.py`](../graph.py) |
| [10](./10-production-roadmap.md) | Scaling, security, retrieval, and storage | Docker and configuration |

## Architecture you will learn

```text
Browser
  | question + thread_id
  v
Flask API
  | root OpenTelemetry span
  v
LangGraph
  retrieve_documentation -> generate_answer -> persist_memory
       |                         |                    |
       +---------- typed AgentState ----------------+
                              |
                  SQLite checkpoint by thread_id

Every request/node -> Tempo spans + Loki events with the same trace_id
Two trace IDs      -> comparison API -> waterfall + correlated logs
```

The workflow is deliberately deterministic: the application chooses the node
order. The model generates an answer but does not control graph routing yet.
Lesson 9 shows where dynamic routing can be added safely.
