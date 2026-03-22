# 4.3 Hands-on with Flink SQL concepts

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/cloud/current/flink/overview.html

The Kafka 101 course uses Flink SQL to demonstrate stream processing. This repo does not currently include a Flink deployment, so the local tutorial focuses on the concepts you would validate.

## Conceptual exercise flow

1. Define a source table over Kafka topics.
2. Write SQL that filters, aggregates, or joins event streams.
3. Emit the result into another topic or sink.
4. Verify continuously updated results as new events arrive.

## Example questions to answer

- How many orders are created per minute?
- What is the running total sales value by region?
- Which users triggered more than three failed login attempts in five minutes?

## Repo gap

To make this executable here, the repo would need:

- a Flink runtime
- SQL scripts
- sample source topics with example data

Prev: [02_stateful_processing.md](02_stateful_processing.md)
