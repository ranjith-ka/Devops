# 5.2 Managed and self-managed models

References:

- https://developer.confluent.io/courses/apache-kafka/events/
- https://docs.confluent.io/platform/current/overview.html

Kafka can be run in multiple ways:

- self-managed on VMs or Kubernetes
- vendor-managed cloud service
- hybrid combinations

## Tradeoffs

| Model | Strength | Cost |
|---|---|---|
| Self-managed | maximum control | more operational burden |
| Managed service | faster setup and operations | less low-level control |

## Where this repo fits

Your repo is on the self-managed side for learning purposes: Strimzi on Minikube, local manifests, and direct Kubernetes ownership.

Prev: [01_kafka_vs_platform.md](01_kafka_vs_platform.md) · Next: [03_what_matters_for_this_repo.md](03_what_matters_for_this_repo.md)
