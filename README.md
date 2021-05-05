# Docker and K8S

Golang, Docker and Kube Practise session

Kubernetes 1.6+

## Helm

```bash
brew install helm
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

## Kind environment

Testing in kind cluster, port mapping required for docker image of Kube worker node. So please make sure extraport mappings are added in the kind/config.yaml
Remember to add in /etc/hosts (to nginx to work)

```bash
CONTAINER ID        IMAGE                   COMMAND                  CREATED             STATUS              PORTS                       NAMES
89c1110261bb        kindest/node:v1.16.15   "/usr/local/bin/entr…"   13 minutes ago      Up 13 minutes       127.0.0.1:65273->6443/tcp   openfaas-control-plane
84a1f8bc9b54        kindest/node:v1.16.15   "/usr/local/bin/entr…"   13 minutes ago      Up 13 minutes       0.0.0.0:32080->32080/tcp    openfaas-worker
```

```bash
$ helm install -f minikube/dev/canary.yaml canary-dev charts/dev
$ helm install -f minikube/dev/prd.yaml prd-dev charts/dev

$ curl -s -H "testing: always" http://st1-dev-vnext.example.com:32080/dev
Welcome to my canary website!%

$ curl -s -H "testing: never" http://st1-dev-vnext.example.com:32080/dev
Welcome to my prod website!%
```

## Install Cobra

### Dadjoke CLI Tool

- Text tutorial: <https://divrhino.com/articles/build-command-line-tool-go-cobra>
- Video tutorial: <https://www.youtube.com/watch?v=-tO7zSv80UY>

Just trying out the tutorial

```bash
cobra init --pkg-name github.com/ranjith-ka/Docker
go mod init github.com/ranjith-ka/Docker
```

Add new command

```bash
cobra add random
```

Used below to convert JSON To go Struct online.

<https://mholt.github.io/json-to-go/>

Added the Pluing REST Client for postman things.

ctrl + alt + M  -- Stop the running code.

### Go Learning

<https://github.com/StephenGrider/GoCasts>
