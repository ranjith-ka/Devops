# 3.3 Hands-on connect workflow

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/platform/current/connect/index.html

This section translates the Kafka 101 Connect exercise into a practical local workflow.

## Suggested learning path

1. Start with the Chapter 1 Kafka cluster.
2. Deploy Kafka Connect as a separate service.
3. Install or include the connector plugin you need.
4. Register a source or sink connector over the Connect REST API.
5. Verify records move between Kafka and the external system.

## What to check

- worker pods are healthy
- connector status is `RUNNING`
- tasks are assigned
- source topics receive records or sink targets receive output

## Repo gap

This repo does not currently include:

- a Kafka Connect deployment
- connector plugin packaging
- example connector manifests

That means this chapter documents the architecture and workflow, but the repo would need extra manifests to run the exercise directly.

Prev: [02_source_and_sink_connectors.md](02_source_and_sink_connectors.md)
