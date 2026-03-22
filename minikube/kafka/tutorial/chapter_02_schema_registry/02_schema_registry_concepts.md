# 2.2 Schema Registry concepts

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/platform/current/schema-registry/fundamentals/index.html

Schema Registry is a service that stores schemas separately from the Kafka broker and lets producers and consumers agree on event structure.

## Main concepts

| Concept | Meaning |
|---|---|
| Subject | A logical name under which schema versions are stored |
| Version | A specific revision of a schema |
| Compatibility | Rules that decide whether a new version is allowed |
| Serializer / deserializer | Client-side code that writes or reads data using the schema |

## Typical flow

1. A producer uses a serializer.
2. The serializer checks whether the schema already exists in Schema Registry.
3. If needed, the schema is registered as a new version.
4. The producer writes the event to Kafka with schema metadata.
5. The consumer fetches the schema and deserializes the record safely.

## Compatibility modes

Common compatibility ideas:

- **backward**: new consumers can read old data
- **forward**: old consumers can read new data
- **full**: both directions are protected

The exact rule choice depends on deployment policy and how much replay compatibility you require.

## Practical takeaway

Schema Registry is not just a serializer helper. It is a governance layer for event contracts.

Prev: [01_why_schemas_matter.md](01_why_schemas_matter.md) · Next: [03_hands_on_schema_workflow.md](03_hands_on_schema_workflow.md)
