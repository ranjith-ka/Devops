# Chapter 1 — Getting Started

References:

- Kafka 101 course overview: https://developer.confluent.io/courses/apache-kafka/events/
- Apache Kafka quick start: https://kafka.apache.org/quickstart

## Sections (this folder)

| # | Topic | Local |
|---|-------|-------|
| 1 | Introduction to Kafka | [01_introduction.md](01_introduction.md) |
| 2 | Topics and partitions | [02_topics_partitions.md](02_topics_partitions.md) |
| 3 | Brokers and replication | [03_brokers_replication.md](03_brokers_replication.md) |
| 4 | Producers and consumers | [04_producers_consumers.md](04_producers_consumers.md) |
| 5 | Running Kafka on Minikube | [05_running_on_minikube.md](05_running_on_minikube.md) |

## Course modules covered

- Introduction
- Topics
- Partitions
- Brokers
- Replication
- Producers
- Hands-on Exercise: Kafka Producer
- Consumers
- Hands-on Exercise: Kafka Consumer

## What this chapter covers

This chapter follows the core sequence from the Kafka 101 playlist:

1. What Kafka is and why event streaming matters.
2. How data is organized into topics and partitions.
3. How brokers store and replicate data.
4. How producers write and consumers read.
5. How to map those concepts to the Strimzi manifests in this repo.

## Kafka in this repo

This repo uses Strimzi to run Kafka on Kubernetes:

- Cluster manifest: [../../kafka.yaml](../../kafka.yaml)
- Kafka UI config: [../../kafka-ui.yaml](../../kafka-ui.yaml)
- Operational notes: [../../Readme.md](../../Readme.md)

Important detail: this manifest uses **ZooKeeper-based** deployment because `zookeeper:` is defined in [../../kafka.yaml](../../kafka.yaml). Newer Kafka deployments may use **KRaft**, but the concepts in this chapter still apply.

## Next in this repo

[Chapter 2 — Schema and event contracts](../chapter_02_schema_registry/README.md)
