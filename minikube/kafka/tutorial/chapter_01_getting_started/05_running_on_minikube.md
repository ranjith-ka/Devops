# 1.5 Running Kafka on Minikube

References:

- Apache Kafka quick start: https://kafka.apache.org/quickstart
- Strimzi install docs: https://strimzi.io/quickstarts/

This section maps Kafka basics to the Kubernetes manifests already present in this repo.

## 1. Create the namespace

```bash
kubectl create namespace kafka
```

## 2. Install the Strimzi operator

```bash
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
```

Wait until the operator pod is ready:

```bash
kubectl get pods -n kafka
```

## 3. Create the Kafka cluster

From repo root:

```bash
kubectl apply -f minikube/kafka/kafka.yaml
```

This creates:

- a 3-broker Kafka cluster
- a 1-node ZooKeeper ensemble
- Strimzi topic and user operators

## 4. Verify the cluster

```bash
kubectl get kafka -n kafka
kubectl get pods -n kafka
kubectl get svc -n kafka
```

The key bootstrap service used by clients is:

- `my-cluster-kafka-bootstrap:9092`

## 5. Produce and consume messages

Open a producer:

```bash
kubectl -n kafka run kafka-producer -ti \
  --image=quay.io/strimzi/kafka:0.39.0-kafka-3.6.1 \
  --rm=true --restart=Never -- \
  bin/kafka-console-producer.sh \
  --bootstrap-server my-cluster-kafka-bootstrap:9092 \
  --topic my-topic
```

Open a consumer in another terminal:

```bash
kubectl -n kafka run kafka-consumer -ti \
  --image=quay.io/strimzi/kafka:0.39.0-kafka-3.6.1 \
  --rm=true --restart=Never -- \
  bin/kafka-console-consumer.sh \
  --bootstrap-server my-cluster-kafka-bootstrap:9092 \
  --topic my-topic \
  --from-beginning
```

Type lines into the producer and confirm they appear in the consumer.

## 6. Optional UI

The file [../../kafka-ui.yaml](../../kafka-ui.yaml) is configured to connect to:

```yaml
bootstrapServers: my-cluster-kafka-bootstrap:9092
```

It is an application config file for Kafka UI, not a full Kafka cluster manifest.

## 7. What to learn next

After this chapter, the next useful topics are:

- message keys and partition strategy
- retention and compaction
- schema management
- Kafka Connect
- stream processing

Prev: [04_producers_consumers.md](04_producers_consumers.md)
