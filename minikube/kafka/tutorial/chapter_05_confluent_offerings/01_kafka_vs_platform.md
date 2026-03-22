# 5.1 Apache Kafka vs platform capabilities

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/platform/current/overview.html

Apache Kafka is the open-source event streaming core. Platforms around Kafka add surrounding capabilities such as:

- schema governance
- managed connectors
- stream processing tooling
- security and RBAC features
- observability and operational tooling

## Practical takeaway

You should separate:

- what Kafka itself provides
- what the surrounding platform adds

That avoids confusion when reading tutorials from different vendors.

Next: [02_managed_and_self_managed.md](02_managed_and_self_managed.md)
