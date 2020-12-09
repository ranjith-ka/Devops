# PSP in minikube 1.13

Start the minikube server with required version

```bash
minikube start --kubernetes-version v1.13.10 --vm-driver=virtualbox --memory=4g --cpus=4
```

```bash
minikube start
kubectl apply -f /path/to/psp.yaml
minikube stop
minikube start --extra-config=apiserver.enable-admission-plugins=PodSecurityPolicy
```

psp.yaml -> Have 2 sections kube-system have privilleged pods
