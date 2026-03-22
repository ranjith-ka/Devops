# 4.1 Stream processing fundamentals

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/cloud/current/flink/overview.html

Stream processing means running computations on events as they arrive.

## Typical operations

- filtering
- transformation
- enrichment
- aggregation
- joins between streams or between streams and tables

## Why Kafka fits this model

Kafka stores ordered event logs that processors can read continuously. That makes it a natural foundation for real-time pipelines.

## Examples

- total orders per minute
- fraud detection on live payments
- user activity sessionization
- inventory updates from purchase events

Next: [02_stateful_processing.md](02_stateful_processing.md)
