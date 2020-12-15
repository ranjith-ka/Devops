# Install ArgoCD

To install ArgoCD in Kind/minikube with Nginx

```bash
brew install kind
kind create cluster --config environment/kind/config.yaml --name agro
helm install nginx stable/nginx-ingress -f minikube/nginx/values.yaml --version 1.36.3
helm install -f minikube/argocd/values.yaml argo  argo/argo-cd
```

## Login to the ArgoCD UI

<http://st2-dev-vnext.example.com:32080>

username: admin
password: 'copy the pod name of argo-server'

Notes:

* Ingress redirection not working as expected, hence creating other URL(any help is appreciated)

## Login to Argo CLI

`argocd login st1-dev-vnext.example.com:32443`

username: admin
password: 'copy the pod name of argo-server'

```bash
argocd app  list
NAME  CLUSTER  NAMESPACE  PROJECT  STATUS  HEALTH  SYNCPOLICY  CONDITIONS  REPO  PATH  TARGET
```
