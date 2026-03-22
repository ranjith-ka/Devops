# 4.2 Stateful processing

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/cloud/current/flink/overview.html

Some stream operations need more than the current event. They need remembered context.

## Examples of state

- current account balance
- events seen in the last 5 minutes
- rolling window counts
- latest profile for a customer

## Time and windows

In streaming systems, time handling matters:

- event time
- processing time
- tumbling windows
- hopping windows
- session windows

## Why this is harder than simple consume-and-print

Stateful processing needs:

- durable state management
- fault recovery
- checkpointing
- careful handling of late or out-of-order events

Prev: [01_stream_processing_fundamentals.md](01_stream_processing_fundamentals.md) · Next: [03_hands_on_flink_sql_concepts.md](03_hands_on_flink_sql_concepts.md)
