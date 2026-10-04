# Flux canary and production lab

This project deploys the existing `charts/dev` chart as `canary-dev` and `prd-dev`.
Source-watcher builds an independent chart artifact for each release, merging
`minikube/dev/canary.yaml` or `minikube/dev/prd.yaml` into that chart's `values.yaml`.

```text
GitRepository/devops → ArtifactGenerator/devops-charts
                       ├─ ExternalArtifact/canary-chart → HelmRelease/canary-app
                       └─ ExternalArtifact/prd-chart    → HelmRelease/prd-app
```

## Run the lab

Use a current Flux CLI (minimum v2.7.2), kind, kubectl and a running Docker engine.
Run commands from the repository root. If you already have a cluster, skip `make kind`.

```bash
make kind
kubectl config current-context
make flux
flux check
kubectl -n flux-system get deployment source-watcher
kubectl get crd artifactgenerators.source.extensions.fluxcd.io
kubectl apply -k minikube/flux/staging
```

The source uses public HTTPS, so no Git secret is needed for this read-only lab.
It reads the remote `main` branch: local chart and values edits only take effect
once committed and pushed there. For a fork, change the URL and branch in
`staging/source.yaml`. For private Git or image automation, configure SSH/token
credentials separately; see tutorial 01 and 05.

Use `-k`, which selects the application resources listed in `staging/kustomization.yaml`.
Notification and image automation examples are optional and are not included.

```bash
flux reconcile source git devops -n default
kubectl -n default wait --for=condition=Ready gitrepository/devops --timeout=2m
kubectl -n default wait --for=condition=Ready artifactgenerator/devops-charts --timeout=2m
kubectl -n default get externalartifacts
flux get helmreleases -n default
kubectl -n default get pods,svc,ingress
```

The ingress controller uses class `nginx` and a NodePort Service. For this local
kind lab it explicitly publishes `127.0.0.1` as the Ingress address; kind maps host
ports 80 and 443 to the controller's NodePorts. This is a reported local address,
not an allocated external load-balancer IP. Allow up to a minute for status updates.

```bash
curl -H 'Host: awesome-http.example.com' http://localhost/dev
```

Both application Ingresses use the same host and `/dev` path so NGINX can attach
the canary backend. Force a specific release using the `testing` header:

```bash
# Canary: Welcome to my website!
curl -H 'Host: awesome-http.example.com' -H 'testing: always' http://localhost/dev
# Production: Welcome to my prod website!
curl -H 'Host: awesome-http.example.com' -H 'testing: never' http://localhost/dev
```

Without that header, the configured canary weight sends approximately 30% of
requests to canary; a small sample will not necessarily match that percentage.

Read [tutorial 06](Tutorials/06_source_watcher.md) for the change-isolation exercise,
troubleshooting and migration details. Earlier tutorials introduce the other controllers
and include historical API examples; the staging manifests use current stable APIs.

## Optional exercises

- [01: Installation](Tutorials/01_Install_Flux.md)
- [02: Sources](Tutorials/02_Source_Controller_flux.md)
- [03: Kustomize and Helm](Tutorials/03_kustomize_helm_controller.md)
- [04: Notifications](Tutorials/04_notification_controller.md): configure your own provider secret before using it.
- [05: Image automation](Tutorials/05_image_automation_ctrlr.md): configure Git write credentials and choose one policy. The old staging example defines the same policy name twice and must be adapted before use.
- [06: Source-watcher](Tutorials/06_source_watcher.md)
