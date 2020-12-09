# Kubernetes in Docker (KIND)

## Pre-requsites

Assuming Docker already running.

```zsh
brew install kind
```

To Install the version, check out the release documentation.

<https://github.com/kubernetes-sigs/kind/releases>

```bash
kind create cluster --config environment/kind/config.yaml --name openfaas
```

config.yaml is for kube 1.16.15 version
