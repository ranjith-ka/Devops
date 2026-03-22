# 2.1 Why schemas matter

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/platform/current/schema-registry/fundamentals/index.html

Kafka can carry bytes, strings, JSON, Avro, Protobuf, and other encodings. The broker itself does not enforce your business structure. That flexibility is useful, but it can become a problem when data producers and consumers evolve independently.

## Typical problems without schemas

- field names drift over time
- consumers disagree about required fields
- type changes break old readers
- teams copy event definitions into multiple services

## What a schema gives you

A schema defines the expected structure of an event:

- field names
- field types
- optional vs required fields
- allowed evolution rules

Example event idea:

```json
{
  "orderId": "ORD-1001",
  "customerId": "CUS-44",
  "amount": 149.95,
  "createdAt": "2026-03-22T10:15:00Z"
}
```

## Why this matters in Kafka

Kafka keeps data for later reprocessing and for multiple consumers. That means a schema change does not only affect today's producer and today's consumer. It affects:

- new consumers
- old consumers
- replay jobs
- downstream pipelines

Next: [02_schema_registry_concepts.md](02_schema_registry_concepts.md)
