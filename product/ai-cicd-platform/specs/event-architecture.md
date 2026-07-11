# Event Architecture

## Event Bus

Use Kafka for durable integration events and NATS for low-latency orchestration and command fanout.

## Event Topics

- `pipeline.run.created`
- `pipeline.run.completed`
- `pipeline.run.failed`
- `deployment.run.created`
- `deployment.run.completed`
- `deployment.run.failed`
- `test.results.ingested`
- `logs.ingested`
- `security.scan.completed`
- `ai.insight.generated`
- `ai.memory.updated`
- `notification.emitted`
- `audit.event.recorded`

## Event Envelope

```json
{
  "event_id": "uuid",
  "event_type": "pipeline.run.completed",
  "schema_version": 1,
  "tenant_id": "org_123",
  "correlation_id": "uuid",
  "causation_id": "uuid",
  "occurred_at": "2026-07-11T00:00:00Z",
  "source": "github-actions",
  "payload": {}
}
```

## Principles

- Idempotent consumers.
- Correlation IDs across API, queue, trace, and audit records.
- Schema versioning for all events.
- Dead-letter queues for failed enrichment or policy validation.
- Event consumers never call model providers on the hot write path.

## AI Event Flow

1. Delivery event lands on Kafka.
2. Ingestion worker enriches metadata.
3. Retrieval index updates OpenSearch and Qdrant.
4. AI summarizer generates an insight or memory update.
5. Notification service fans out the result.
