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

```
minikube start --docker-env HTTP_PROXY=http://192.168.1.9:10809 --docker-env HTTPS_PROXY=http://192.168.1.9:10809 --docker-env NO_PROXY=192.168.99.0/24
export no_proxy=$no_proxy,$(minikube ip)
```
