# ES in minikube

```bash
helm install -f minikube/elastic-search/master.yaml elasticsearch-master elastic/elasticsearch --version 7.9.2
helm install -f minikube/elastic-search/data.yaml elasticsearch-data elastic/elasticsearch --version 7.9.2
helm install -f minikube/elastic-search/client.yaml elasticsearch-client elastic/elasticsearch --version 7.9.2
```

```bash
minikube start --kubernetes-version v1.16.8 --vm-driver=virtualbox --memory=4g --cpus=2
```
