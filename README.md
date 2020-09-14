# Docker and K8S

Docker and Kube Practise session

Kubernetes 1.6+

## Helm

Add the stable repo for some default service testing.

```bash
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
```

## Automated PR

```bash
brew install github/gh/gh
git add .
git commit -am "just testing"
gh pr create -f
```

## Create Nginx Service

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm repo update
```

<https://github.com/kubernetes/ingress-nginx/tree/master/charts/ingress-nginx>

```bash
helm install -f minikube/nginx/values.yaml nginx ingress-nginx/ingress-nginx
```

```bash
$ minikube service ingress-nginx-controller  --url
http://192.168.99.100:32080
http://192.168.99.100:31443
http://192.168.99.100:32443
```

Add st1-dev-vnext.example.com in /etc/hosts to connect local

```bash
curl http://st1-dev-vnext.example.com:32080/dev
```

```bash
$ helm install -f minikube/dev/canary.yaml canary-dev charts/dev
$ helm install -f minikube/dev/prd.yaml prd-dev charts/dev

$ curl -s -H "testing: always" http://st1-dev-vnext.example.com:32080/dev
Welcome to my canary website!%

$ curl -s -H "testing: never" http://st1-dev-vnext.example.com:32080/dev
Welcome to my prod website!%
```
