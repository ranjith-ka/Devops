### Jager

helm install jaeger jaegertracing/jaeger \
 --set provisionDataStore.cassandra=false \
 --set storage.type=elasticsearch \
 --set storage.elasticsearch.host=elasticsearch-master \
 --set storage.elasticsearch.port=9200
