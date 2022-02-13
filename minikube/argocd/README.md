# Install ArgoCD

To install ArgoCD in Kind/minikube with Nginx

```bash
brew install kind
kind create cluster --config environment/kind/config.yaml --name agro
helm install -f minikube/nginx/values.yaml nginx ingress-nginx/ingress-nginx --version 3.36.0
helm install -f minikube/argocd/values.yaml argo  argo/argo-cd
```

Just to update the version tested here in this demo.

```bash
✗ helm ls
NAME  NAMESPACE REVISION UPDATED                              STATUS   CHART                APP VERSION
argo  default   1        2021-09-07 21:03:59.591332 +0530 IST deployed argo-cd-3.17.5       2.1.1
nginx default   1        2021-09-07 21:02:50.596626 +0530 IST deployed ingress-nginx-3.36.0 0.49.0
```

## Login to the ArgoCD UI

<http://awesome-tcp.example.com>

username: admin
password: 'get password'

example:
'kubectl get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d'

Notes:

-   Ingress redirection not working as expected, hence creating other URL(any help is appreciated)

## Login to Argo CLI

`argocd login awesome-http.example.com`

username: admin
password: 'copy the pod name of argo-server'

```bash
argocd app  list
NAME  CLUSTER  NAMESPACE  PROJECT  STATUS  HEALTH  SYNCPOLICY  CONDITIONS  REPO  PATH  TARGET
```
