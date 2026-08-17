# Lesson 6: Tempo node tracing

## Learning objective

Understand parent-child span creation and make every future graph node visible in
the waterfall.

## 1. Trace setup

`tracing.py` creates an OpenTelemetry `TracerProvider`, resource attributes, an
OTLP/gRPC exporter, and a batch span processor. The application exports to the
Collector, which forwards spans to Tempo.

## 2. Node instrumentation

Every graph node is registered through `traced_node(name, function)`. The wrapper:

1. starts `graph.node.<name>`;
2. records the stable `langgraph.node` attribute;
3. emits start/completion Loki events;
4. records exceptions and error logs;
5. re-raises failures so graph semantics remain correct.

Because graph execution occurs inside `langchain.request`, OpenTelemetry creates
the correct parent-child relationship automatically.

## 3. Attribute design

Good attributes are bounded and useful for grouping: node name, model name,
environment, result count, and status. Avoid prompts, document bodies, raw user
IDs, and arbitrary tool results.

## Exercise

Add a 100 ms delay to one node, generate before/after traces, and compare them in
the application. Remove the delay afterward. Confirm the difference is assigned
to the expected node.

## Common mistakes

- Starting node work before entering the span context.
- Swallowing exceptions in the instrumentation wrapper.
- Shutting down the provider after every Flask request.
- Creating high-cardinality resource labels.

## Next lesson

Lesson 7 adds the event detail that spans intentionally do not contain.
