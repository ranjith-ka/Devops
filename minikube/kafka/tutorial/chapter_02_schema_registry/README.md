# Chapter 2 — Schema and Event Contracts

References:

- Confluent Schema Registry module: https://developer.confluent.io/courses/apache-kafka/events/
- Schema Registry fundamentals: https://docs.confluent.io/platform/current/schema-registry/fundamentals/index.html

## Sections (this folder)

| # | Topic | Local |
|---|-------|-------|
| 1 | Why schemas matter | [01_why_schemas_matter.md](01_why_schemas_matter.md) |
| 2 | Schema Registry concepts | [02_schema_registry_concepts.md](02_schema_registry_concepts.md) |
| 3 | Hands-on schema workflow | [03_hands_on_schema_workflow.md](03_hands_on_schema_workflow.md) |

## Course modules covered

- Confluent Schema Registry
- Hands-on Exercise: Confluent Schema Registry

## What this chapter covers

This chapter explains how Kafka events evolve safely over time:

1. Why raw JSON without contracts becomes hard to manage.
2. How Schema Registry stores and versions schemas.
3. How compatibility rules protect producers and consumers.
4. How to approach the hands-on exercise in a Kubernetes environment.

## Important repo-specific note

This repo currently contains a Kafka cluster manifest and Kafka UI configuration, but it does **not** yet include a deployed Schema Registry manifest. The hands-on section therefore focuses on the workflow and the components you would add next.

Prev: [../chapter_01_getting_started/README.md](../chapter_01_getting_started/README.md) · Next: [../chapter_03_kafka_connect/README.md](../chapter_03_kafka_connect/README.md)
