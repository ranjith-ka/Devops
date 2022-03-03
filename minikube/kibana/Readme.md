# Install Kibana in Minikube

```bash
helm install -f minikube/kibana/values.yaml kibana elastic/kibana
```

If you try in minikube, sync time is required.

```bash
minikube ssh -- docker run -i --rm --privileged --pid=host debian nsenter -t 1 -m -u -n -i date -u $(date -u +%m%d%H%M%Y)
```
