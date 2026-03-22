# 2.3 Hands-on schema workflow

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/platform/current/schema-registry/fundamentals/index.html

This is the repo-oriented equivalent of the Kafka 101 Schema Registry exercise.

## Goal

Move from plain unstructured messages toward versioned event contracts.

## Suggested local flow

1. Keep the Kafka cluster from Chapter 1 running.
2. Add or deploy a Schema Registry service.
3. Create a test topic such as `orders`.
4. Produce events using a schema-aware client.
5. Evolve the schema by adding an optional field.
6. Verify that consumers still read historical and new records correctly.

## Example schema evolution

Version 1 idea:

```json
{
  "type": "record",
  "name": "OrderCreated",
  "fields": [
    {"name": "orderId", "type": "string"},
    {"name": "amount", "type": "double"}
  ]
}
```

Version 2 idea with backward-compatible change:

```json
{
  "type": "record",
  "name": "OrderCreated",
  "fields": [
    {"name": "orderId", "type": "string"},
    {"name": "amount", "type": "double"},
    {"name": "currency", "type": ["null", "string"], "default": null}
  ]
}
```

## What to verify

- the new schema version is accepted
- the producer can write the new event shape
- old data remains readable
- the consumer handles the new optional field correctly

## Repo gap

If you want this tutorial to become executable end-to-end inside this repo, the next step is to add:

- a Schema Registry deployment manifest
- an example producer and consumer using Avro or Protobuf

Prev: [02_schema_registry_concepts.md](02_schema_registry_concepts.md)
