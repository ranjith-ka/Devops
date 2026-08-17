# Lesson 7: Loki log correlation

## Learning objective

Understand structured event design, trace correlation, and why logs are
best-effort in this application.

## 1. Log structure

`emit_log()` creates JSON containing timestamp, level, event, trace ID, and safe
event fields. Loki stream labels contain only service and environment.

The trace ID stays in the JSON body. Making every trace ID a Loki label would
create unbounded label cardinality.

## 2. Current events

- `request.started` / `request.completed`
- `node.started` / `node.completed` / `node.failed`
- `documentation.retrieved`
- `model.completed`
- `memory.persisted`

These events explain what occurred without storing prompts, answers, or document
bodies.

## 3. Failure behavior

Loki calls have short timeouts and failures are ignored. An observability outage
must not turn into an application outage. A production deployment should send
logs to stdout and use an agent such as Grafana Alloy for buffering and shipping.

## Exercise

Search Loki for one trace ID, then stop Loki and submit another question. Confirm
that the answer still succeeds. Restart Loki before continuing.

## Security checkpoint

Before adding a field, ask:

- Could it contain credentials or personal data?
- Is it bounded in size?
- Is it needed for diagnosis?
- Should it be a body field rather than a label?

## Next lesson

Lesson 8 joins Tempo timing and Loki events for two real executions.
