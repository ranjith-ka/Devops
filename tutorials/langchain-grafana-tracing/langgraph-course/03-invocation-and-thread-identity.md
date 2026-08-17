# Lesson 3: Invocation and thread identity

## Learning objective

Understand request validation, root spans, graph invocation, and the difference
between thread IDs and trace IDs.

## 1. Two different identifiers

| Identifier | Lifetime | Purpose |
| --- | --- | --- |
| `thread_id` | Many user turns | Selects LangGraph checkpoint history |
| `trace_id` | One request | Correlates spans and logs for one execution |

Reusing a trace ID would incorrectly combine separate requests. Generating a new
thread ID for every question would erase conversational continuity.

## 2. Browser behavior

The browser stores `langgraph-thread-id` in `localStorage`. Every question from
that browser profile reuses it until storage is cleared.

The API limits the incoming thread ID length before placing it in span metadata.
In a multi-user system, the server should derive this value from an authenticated
session instead of trusting the browser.

## 3. Root-span boundary

`ui.py` starts `langchain.request` before invoking the graph. Node spans created
inside that context automatically become children of the root span.

The API returns the root trace ID so the user can query Tempo and Loki.

## Exercise

Open the application in a normal window and a private window. Ask two follow-up
questions in each. Confirm that memory continues within one window but does not
cross between the two thread IDs.

## Extension checkpoint

Before production, replace browser-generated identity with:

```text
authenticated user + explicit conversation ID + tenant boundary
```

Never use raw email addresses or other personal identifiers as indexed telemetry
labels.

## Next lesson

Lesson 4 explains exactly how SQLite checkpoint persistence survives container
rebuilds.
