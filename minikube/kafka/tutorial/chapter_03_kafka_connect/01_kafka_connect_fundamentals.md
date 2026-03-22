# 3.1 Kafka Connect fundamentals

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/platform/current/connect/index.html

Kafka Connect is a framework for moving data into and out of Kafka using reusable connectors.

## Why it exists

Without Kafka Connect, teams often write custom code for:

- database ingestion
- object storage export
- search indexing
- SaaS integration

Kafka Connect standardizes that work.

## Main building blocks

| Component | Role |
|---|---|
| Connect worker | Runs the Connect runtime |
| Connector | Defines the integration type and configuration |
| Task | Unit of parallel work created by a connector |

## Distributed mindset

Kafka Connect is usually run as a service cluster, not as a one-off script. That allows:

- central configuration
- scaling tasks horizontally
- restart and fault tolerance
- shared management APIs

Next: [02_source_and_sink_connectors.md](02_source_and_sink_connectors.md)
