# Kubernetes in Docker (KIND)

config.yaml is for kube 1.16.15 version

<https://sookocheff.com/post/kubernetes/local-kubernetes-development-with-kind/>
<https://kind.sigs.k8s.io/docs/user/local-registry/> - Hope this works after years.

## Pre-requsites

Assuming Docker already running.

```zsh
brew install kind helm
```

To Install the version, check out the release documentation.

<https://github.com/kubernetes-sigs/kind/releases>

```bash
kind create cluster --config kind/config.yaml --name k8s
```

If we need a local registry, then connect the docker host with Kind

```bash
docker network connect "kind" "kind-registry"
```

To verify the image from registry, some time we need to list out the image in local registry.

```bash
docker tag myimage:tag localhost:5000/myimage:tag
docker push localhost:5000/myimage:tag
curl -X GET http://localhost:5000/v2/_catalog
curl -X GET http://localhost:5000/v2/myimage/tags/list
```

Note:
Running registry means we need to build local and push to localhost:5000
Follow the above documentation.

## Makefile

Just run the make command to run the k8s and install nginx and sample application.

```bash
$ make build-image
$ make kind-cluster
$ make load-image
$ make ingress
$ make install-app
```

```bash
$ curl -s -H "testing: always" http://st1-dev-vnext.example.com/dev
Welcome to my canary website!%
$ curl -s -H "testing: never" http://st1-dev-vnext.example.com/dev
Welcome to my prod website!%
```
