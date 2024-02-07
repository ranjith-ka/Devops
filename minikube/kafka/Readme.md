# Kafka

## Install steps

`kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka`


### Prodcuer

`kubectl -n kafka run kafka-producer -ti --image=quay.io/strimzi/kafka:0.39.0-kafka-3.6.1 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic`


### Consume the Topic

`kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.39.0-kafka-3.6.1 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning`


### Features

Kafka UI
ACL
SSL & TLS
Monitoring
Autoscaling - Broker & Topics
Mirror making
Bridge configuration
Schema Registry
kafka connect
