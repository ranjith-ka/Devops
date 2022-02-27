### Minikube testing

Minikube can be used as docker desktop replacement

```bash
$ make minikube
$ eval $(minikube -p minikube docker-env)
$ make snapshot
$ make ingress
$ make install-app
$ make monitoring
```

Wait for the pod to come online.

```bash
$ kubectl get po --watch
```

```bash
minikube delete
```
