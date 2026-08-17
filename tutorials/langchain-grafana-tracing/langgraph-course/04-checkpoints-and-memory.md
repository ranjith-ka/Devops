# Lesson 4: Checkpoints and durable memory

## Learning objective

Understand what LangGraph checkpoints, how `thread_id` restores it, and how this
differs from documentation retrieval.

## 1. Checkpointer setup

`AgentGraph` opens a SQLite connection and compiles the graph with a
`SqliteSaver`. Each invocation includes:

```python
config = {"configurable": {"thread_id": thread_id}}
graph.invoke({"question": question}, config=config)
```

LangGraph loads the latest state for that thread, applies the new question, runs
the graph, and writes new checkpoints.

## 2. What memory contains

The `persist_memory` node appends user and assistant messages and keeps the most
recent 20 entries. The generation node uses the latest six entries to bound the
prompt size.

These are two separate limits:

- stored history limit controls checkpoint growth;
- prompt history limit controls model context and latency.

## 3. Docker persistence

`CHECKPOINT_DB=/data/checkpoints.sqlite` is stored in the named `app-data`
volume. `docker compose down` preserves it. `docker compose down -v` deletes it.

## 4. Memory is not knowledge retrieval

Conversation memory remembers turns in one thread. Documentation retrieval finds
reference material shared across threads. They have different lifecycle,
security, deletion, and relevance rules.

## Exercise

Ask “Remember that my test color is amber,” restart only the UI container, then
ask for the test color using the same browser. Inspect `/data/checkpoints.sqlite`
inside the container and confirm the volume remains mounted.

## Production extension

SQLite is appropriate for this single UI process. Multiple replicas should use a
shared durable checkpointer such as Postgres, with tenant-aware thread IDs and
retention/deletion policies.

## Next lesson

Lesson 5 explains how Markdown documentation is chunked, ranked, and injected
into the answer node.
