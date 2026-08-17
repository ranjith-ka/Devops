# Lesson 1: Architecture tour

## Learning objective

Understand every major component and follow one request from the browser to its
answer, checkpoint, trace, and logs.

## 1. Runtime components

The Docker stack contains five services:

| Service | Purpose | Port |
| --- | --- | --- |
| UI | Flask API and LangGraph execution | 5001 |
| OpenTelemetry Collector | Receives and forwards spans | 4317/4318 |
| Tempo | Stores traces | 3200 |
| Loki | Stores structured logs | 3100 |
| Grafana | Explores Tempo and Loki | 3000 |

Ollama runs on the host and is reached from Docker through
`host.docker.internal:11434`.

## 2. One request lifecycle

1. The browser creates or restores a `thread_id` from `localStorage`.
2. It sends `question` and `thread_id` to `POST /api/question`.
3. `ui.py` opens the root `langchain.request` span.
4. `AgentGraph.invoke()` runs the compiled graph.
5. The documentation node finds relevant Markdown chunks.
6. The model node receives history, documentation, and the new question.
7. The memory node appends the question and answer to history.
8. LangGraph checkpoints the resulting state in SQLite.
9. The API flushes spans and returns the answer, trace ID, thread ID, and sources.

## 3. Data and observability paths

Do not confuse these paths:

```text
Business data: question -> graph state -> answer
Memory:        graph state -> SQLite checkpoint
Tracing:       spans -> Collector -> Tempo
Logging:       JSON events -> Loki
```

Tempo answers “where was time spent?” Loki answers “what happened at that
moment?” SQLite answers “what did this thread previously know?”

## Exercise

Run one question, copy its trace ID, and locate it in both:

- Grafana Explore with the Tempo datasource;
- Grafana Explore with Loki query `{service_name="langchain-tracing"}`.

Then identify the root request and the three graph nodes.

## Knowledge check

1. Which component owns conversation memory?
2. Why should Loki failure not prevent an answer?
3. What connects a Tempo trace to Loki records?

## Next lesson

Lesson 2 opens the graph itself and explains state, partial updates, nodes, and
edges.
