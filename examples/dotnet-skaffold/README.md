# .NET + Skaffold + Kind

This example builds an ASP.NET Core API with Docker BuildKit, loads the image
into a local Kind cluster, deploys it with Skaffold, and forwards the service to
`http://localhost:8080`.

## Run locally without Kubernetes

```bash
dotnet run
```

## Run with Kind and Skaffold

Create the cluster once:

```bash
kind create cluster --name k8s
```

Verify the context before deploying:

```bash
kubectl config current-context
```

It must print `kind-k8s`. Then run:

```bash
skaffold dev --kube-context kind-k8s
```

Test the API from another terminal:

```bash
curl http://localhost:8080/
curl http://localhost:8080/healthz
```

Press Ctrl+C in the Skaffold terminal to stop and clean up the deployment.

## Build inside Kubernetes with rootless BuildKit

This alternative runs the Dockerfile build in an ephemeral BuildKit pod and
pushes the resulting image to Docker Hub. It is the maintained replacement for
the old in-cluster Kaniko workflow.

Authenticate to Docker Hub first. The account must be allowed to push
`ranjithka/dotnet-skaffold-api`:

```bash
docker login
```

Create a local Kind cluster and the builder namespace:

```bash
kind create cluster --name k8s
kubectl --context kind-k8s create namespace buildkit
kubectl --context kind-k8s --namespace buildkit create secret generic \
  buildkit-registry-auth \
  --from-file=config.json="$HOME/.docker/config.json"
```

Verify that `kind-k8s` is the target, then start the remote build workflow:

```bash
skaffold dev \
  --filename skaffold-remote.yaml \
  --kube-context kind-k8s
```

The registry-auth secret is persistent and reused by every build. The custom
builder creates an ephemeral rootless BuildKit pod, uploads the source, builds
and pushes the image, deletes the builder pod, and deploys the API.

BuildKit stores its local content and layer cache in the persistent
`buildkit-cache` claim. It also imports and exports a registry cache at
`ranjithka/dotnet-skaffold-api:buildcache`, allowing cache reuse if a build is
scheduled on another node. The first build populates both caches; later builds
should show `CACHED` for unchanged Dockerfile steps.

Successful builds delete their ephemeral builder pod. Failed builds preserve it
for inspection; the failure output prints the exact `kubectl delete pod`
command to use after debugging.

This proof-of-concept intentionally pins both the builder and deployer to the
`kind-k8s` context. For another cluster, change `KUBECONTEXT` and
`deploy.kubeContext` in `skaffold-remote.yaml` together.

The builder context can instead be overridden for a remote Kind cluster by
setting `BUILD_KUBE_CONTEXT`; the Skaffold `--kube-context` flag overrides the
deployer context:

```bash
BUILD_KUBE_CONTEXT=kind-remote skaffold dev \
  --filename skaffold-remote.yaml \
  --kube-context kind-remote
```

The cluster must permit an `Unconfined` seccomp profile. Some managed clusters
block it through Pod Security Admission or policy engines; coordinate a
dedicated builder namespace policy before using this configuration there.
