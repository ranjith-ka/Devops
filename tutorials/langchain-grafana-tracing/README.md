# LangChain tracing with Grafana Tempo

This tutorial builds a stateful LangGraph workflow, retrieves local documentation,
stores conversation checkpoints in SQLite, sends OpenTelemetry spans to Tempo,
and sends trace-correlated structured logs to Loki. Grafana and the included UI
can inspect and compare both signals.

The sample uses a deterministic fake model, so the entire lab works without an
LLM account or API key. Replace that one function with a provider model after the
trace path works.

## Step-by-step LangChain Agents course

For the complete agents curriculum, start with
[agents-course/README.md](./agents-course/README.md). Lessons are intentionally
small and should be completed in order.

## Step-by-step course for this LangGraph application

To learn the implemented graph, durable memory, documentation retrieval, Tempo,
Loki, comparison UI, and extension patterns, start with
[langgraph-course/README.md](./langgraph-course/README.md).

## What you will learn

- compose a prompt, model step, and parser with LangChain runnables;
- create parent and child spans around an LLM request;
- attach model, token, latency, and error attributes without leaking prompts;
- inspect a trace in Grafana;
- move from a local Tempo process to a durable, horizontally scalable design;
- control ingestion, storage, and LLM cost.

## Architecture

```text
Python/LangGraph
  request span
    ├── graph.node.retrieve_documentation
    ├── graph.node.generate_answer
    └── graph.node.persist_memory
          │ OTLP/gRPC
          v
OpenTelemetry Collector ──OTLP──> Tempo ──query──> Grafana
Application logs ────────────────> Loki ───query──> Grafana/UI
  memory limit + batches            local disk       Explore
```

The application talks to the Collector, not directly to Tempo. That boundary is
important: sampling, redaction, retries, routing, authentication, and backend
changes can then happen without changing application code.

## 1. Run the lab

Prerequisites: Docker Compose and Python 3.10 or newer.

```bash
cd tutorials/langchain-grafana-tracing
docker compose up -d
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
set -a; source .env.example; set +a
python app.py "How do chains and runnables work?"
```

The command prints a result and a 32-character trace ID. Open
`http://localhost:3000`, choose **Explore**, select **Tempo**, and either:

1. paste the trace ID into the trace lookup field; or
2. use the TraceQL query `{ resource.service.name = "langchain-tutorial" }`.

Allow a few seconds for the SDK batch exporter, Collector, and Tempo to flush.
Inspect the waterfall and confirm the request has documentation retrieval,
answer generation, and memory persistence children.

Useful diagnostics:

```bash
docker compose ps
docker compose logs otel-collector tempo
docker compose down
```

`docker compose down` preserves the named Tempo volume. Use
`docker compose down -v` only when you intentionally want to delete local traces.

## 2. Understand the LangGraph workflow

The graph in `graph.py` uses a typed state and explicit nodes:

```text
START -> retrieve_documentation -> generate_answer -> persist_memory -> END
```

Each node is wrapped in a stable OpenTelemetry span and emits structured Loki
events carrying the same trace ID. SQLite checkpoints use `thread_id` to retain
the last conversation messages across requests and container restarts.

To add a node, define a function returning a partial `AgentState`, register it
with `builder.add_node`, and connect it using `builder.add_edge` or a conditional
edge. Wrap it with `traced_node` so Tempo and Loki comparisons include it.

Markdown files under `docs/` are loaded by `DocumentationStore`. The current
retriever is intentionally local and lexical; it can later be replaced by an
embedding/vector implementation without changing the graph state contract.

The fake model also records:

- `gen_ai.request.model`;
- estimated input and output token counts;
- a prompt hash by default;
- an opt-in prompt preview when `TRACE_CONTENT=true`.

The token counts in this fake model are word counts for teaching purposes. In a
real integration, record the provider's returned usage metadata; do not estimate
billing from `split()`.

## 3. Use a real model

Install the provider package separately, for example `langchain-openai`, then
replace `fake_model` with a configured chat model's `invoke` method:

```python
from langchain.chat_models import init_chat_model

model = init_chat_model(
    "openai:YOUR_MODEL_NAME",
    temperature=0,
    timeout=30,
    max_retries=2,
)

chain = (
    traced_step("prompt.render", "prompt", prompt.invoke)
    | traced_step("model.generate", "chat", model.invoke)
    | traced_step("output.parse", "parse", StrOutputParser().invoke)
)
```

Set the provider credential in the environment, never in source control. Capture
the real response's usage metadata in the model span. Keep provider/model names,
status, duration, and token counts; avoid raw prompts and completions unless a
reviewed debugging policy explicitly allows them.

## 4. Span design that remains useful at scale

Use a stable, low-cardinality vocabulary:

| Level | Span | Keep as attributes |
| --- | --- | --- |
| Request | `langchain.request` | route, tenant tier, status, release |
| Workflow | agent/chain/retrieval | operation, result count, cache hit |
| Model | `model.generate` | provider, model, input/output tokens, retry count |
| Tool | tool name | tool type, status, duration; not arbitrary arguments |

Do not put user IDs, full URLs with query strings, prompts, documents, tool
results, or session IDs into indexed attributes. High-cardinality attributes
increase index and query cost; sensitive content increases breach impact.

Propagate W3C trace context across HTTP, queues, and worker boundaries. A single
trace should connect the incoming API request, retrieval, model calls, tools, and
outgoing dependencies. Use metrics—not trace scans—for alerting on request rate,
error rate, latency percentiles, and token spend.

## 5. Scale the architecture in stages

### Stage A: developer and small workload

Use this repository's single Collector and monolithic Tempo with local disk.
Retention is 24 hours. This is cheap and simple, but it has no high availability
and local storage is not a production durability boundary.

### Stage B: production baseline

Deploy two Collector layers:

```text
apps -> per-node/sidecar collectors -> load balancer -> gateway collectors
                                                   -> managed Tempo or Tempo
```

- Agents handle local batching and shield apps from backend interruptions.
- Stateless gateways perform central redaction, tenant routing, rate limiting,
  and baseline probabilistic sampling.
- Run at least two gateways across failure domains with bounded queues and memory.
- Send to a managed OTLP backend when operating Tempo is not a core capability.
- If self-hosting, use object storage and tested retention/lifecycle policies.

Scale collectors on CPU, memory, refused spans, queue utilization, and export
latency—not only pod count. Apply backpressure and cap queues so an observability
incident cannot exhaust application nodes.

### Stage C: high volume and independent scaling

Use Tempo microservices mode when the write path, recent queries, long-term
queries, and compaction need to scale separately. Tempo 3 microservices mode uses
a Kafka-compatible durable queue and object storage. Scale distributors for
ingestion, live stores for recent searches, queriers/query frontends for reads,
and block builders/workers for storage maintenance.

Partition tenants, set per-tenant ingestion and query limits, isolate noisy
tenants, and test failure of collectors, the queue, object storage, and a full
availability zone. Keep the local Compose stack for development; do not stretch
it into production.

## 6. Cost controls in priority order

1. **Minimize payload first.** Drop prompt/completion bodies and oversized tool
   results. Keep hashes, byte counts, token counts, status, and safe categories.
2. **Head-sample ordinary successes.** Start with 5–10% after measuring traffic.
   Always keep a small unbiased baseline so percentiles and rare behavior remain
   observable.
3. **Tail-sample valuable traces.** At a gateway tier, retain errors, timeouts,
   high latency, high-token calls, selected releases, and a random baseline.
4. **Use tiered retention.** Keep searchable hot traces briefly; use object-store
   lifecycle rules for longer retention only where incident or compliance needs
   justify it.
5. **Turn traces into metrics.** Service graphs and RED metrics are cheaper for
   dashboards and alerts than repeatedly searching raw traces.
6. **Control the LLM itself.** Cap context and output tokens, cache safe repeated
   results, route simple work to smaller models, bound retries, and use budgets per
   tenant.

Estimate daily trace volume before choosing infrastructure:

```text
daily_bytes = requests_per_second
            × spans_per_request
            × average_encoded_span_bytes
            × 86,400
            × retained_sample_rate
```

Example: `100 RPS × 8 spans × 1,000 bytes × 86,400 × 0.10` is about 6.9 GB/day
before compression and backend overhead. Measure actual OTLP traffic because LLM
payload attributes can make the average span far larger than 1 KB.

Track a separate business cost estimate:

```text
llm_cost = input_tokens × input_price + output_tokens × output_price
```

Aggregate it as metrics by service, tenant tier, use case, and model. Do not make
raw tenant identifiers metric labels.

## 7. Production readiness checklist

- TLS and authentication on every OTLP hop; ports are not publicly exposed.
- Prompt/completion capture is off by default and redaction is tested.
- Collector memory limits, queues, retry bounds, and dropped-span alerts exist.
- Trace and tenant IDs are present in structured logs, without duplicating bodies.
- Sampling rules retain errors and a statistically useful baseline.
- Object-store encryption, retention, deletion, backup, and access controls exist.
- Load tests validate peak spans/second and worst-case span size.
- Dashboards cover ingest rate, rejected spans, exporter failures, queue fill,
  query latency, storage growth, token use, and estimated LLM cost.
- Runbooks explain how to reduce sampling or disable content capture during an
  incident without redeploying every application.

## References

- [LangChain overview](https://docs.langchain.com/oss/python/langchain/overview)
- [LangChain installation](https://docs.langchain.com/oss/python/langchain/install)
- [OpenTelemetry Python instrumentation](https://opentelemetry.io/docs/languages/python/instrumentation/)
- [OpenTelemetry Python exporters](https://opentelemetry.io/docs/languages/python/exporters/)
- [Grafana Tempo quick start](https://grafana.com/docs/tempo/latest/docker-example/)
- [Grafana Tempo architecture](https://grafana.com/docs/tempo/latest/introduction/architecture/)
- [Plan a Tempo deployment](https://grafana.com/docs/tempo/latest/set-up-for-tracing/setup-tempo/plan/)


## Documentation

- [LangGraph Architecture](./docs/langgraph-architecture.md)
