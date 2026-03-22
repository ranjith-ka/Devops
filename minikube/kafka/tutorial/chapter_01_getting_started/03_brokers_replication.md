# 1.3 Brokers and replication

Reference: https://developer.confluent.io/courses/apache-kafka/events/

## Broker

A **broker** is a Kafka server. A Kafka cluster is made of one or more brokers.

In this repo, the cluster is configured with:

```yaml
kafka:
  replicas: 3
```

That means Strimzi will create **3 Kafka broker pods**.

## Partition leaders and followers

Each partition has a **leader** broker and, when replication is enabled, one or more **followers**.

- Producers write to the leader.
- Consumers read from the leader.
- Followers replicate the leader's data.

## Replication

Replication protects data when a broker fails.

Common terms:

- **replication factor**: how many copies of partition data exist
- **ISR**: in-sync replicas that are caught up closely enough to the leader
- **leader election**: choosing a new leader after failure

## Important repo-specific detail

The sample manifest is intentionally lightweight for Minikube:

```yaml
offsets.topic.replication.factor: 1
transaction.state.log.replication.factor: 1
default.replication.factor: 1
min.insync.replicas: 1
```

This is fine for local learning, but it is **not** a production durability setup. In production, replication factor and `min.insync.replicas` are typically greater than `1`.

## ZooKeeper note

This manifest also contains:

```yaml
zookeeper:
  replicas: 1
```

So this deployment uses the older ZooKeeper-based metadata architecture. Many modern Kafka installations use **KRaft**, but the topic, partition, producer, and consumer model stays the same.

Prev: [02_topics_partitions.md](02_topics_partitions.md) · Next: [04_producers_consumers.md](04_producers_consumers.md)
