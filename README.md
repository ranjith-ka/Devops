# Docker and K8S

Docker and Kube Practise session

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
