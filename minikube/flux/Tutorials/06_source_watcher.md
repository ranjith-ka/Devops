# Source-watcher with the existing Devops applications

## What this exercise teaches

Source-controller fetches the Git repository. Source-watcher watches its artifact,
extracts the shared Helm chart and merges the release-specific values. It publishes
`canary-chart` and `prd-chart` ExternalArtifacts. Helm-controller consumes those
artifacts through each HelmRelease's `spec.chartRef`.

The copies create a `dev/` chart directory with Chart.yaml, templates and merged
values.yaml. No chart version bump is needed to observe content changes through
ExternalArtifact revisions. All resources live in `default`; Flux controllers live
in `flux-system`.

References: [ArtifactGenerator](https://fluxcd.io/flux/components/source/artifactgenerators/),
[Helm chart references](https://fluxcd.io/flux/components/helm/helmreleases/#chart-reference).

## 1. Install and inspect

Follow the [lab README](../Readme.md) to run `make flux` and apply the staging
Kustomization. The Flux CLI must be v2.7.2 or later; use a current stable release
for the current Kubernetes cluster.

```bash
flux --version
flux check
kubectl -n flux-system get deployments
kubectl -n flux-system get deployment helm-controller \
  -o jsonpath='{.spec.template.spec.containers[0].args}'
```

Helm-controller must enable `ExternalArtifact`. The current Flux installation
includes `--feature-gates=ExternalArtifact=true`. If yours does not, enable it in
your managed controller configuration before using these HelmReleases.

## 2. Inspect the artifacts and releases

```bash
kubectl -n default get gitrepository devops
kubectl -n default get artifactgenerator devops-charts
kubectl -n default get externalartifacts
kubectl -n default get externalartifact canary-chart prd-chart \
  -o custom-columns=NAME:.metadata.name,REVISION:.status.artifact.revision
flux get helmreleases -n default
helm list -n default
```

The generated charts merge base values with the environment values in that order.
The canary HelmRelease also retains its existing inline `spec.values`; those override
the generated values. In particular, its inline image tag remains `0.0.1`.
To practise image changes through the values file, remove the inline image override
first. For the isolation exercise below, use `replicaCount`, which has no inline override.

## 3. Practise independent updates

Record both ExternalArtifact revisions using the command above. Add or change
`replicaCount: 2` in `minikube/dev/canary.yaml`. Commit and push that edit to the
branch watched by `GitRepository/devops`, then run:

```bash
flux reconcile source git devops -n default
kubectl -n default get externalartifact canary-chart prd-chart \
  -o custom-columns=NAME:.metadata.name,REVISION:.status.artifact.revision
flux get helmreleases -n default
kubectl -n default get deployments
```

Allow the controllers time to reconcile. The canary artifact revision should change;
the production artifact revision should stay the same. Production may still perform
its periodic reconciliation, but it should not upgrade because of this canary edit.

Repeat with `minikube/dev/prd.yaml`; only the production artifact should change.
Finally, make a meaningful shared chart change (for example a default value used by
both releases in `charts/dev/values.yaml`): both artifacts should change.
A documentation-only edit outside the copied chart and values files should change
neither generated artifact. Restore the exercise edits through Git when finished.

## 4. Troubleshoot each stage

```bash
kubectl -n default describe gitrepository devops
kubectl -n default describe artifactgenerator devops-charts
kubectl -n default describe externalartifact canary-chart
kubectl -n default describe helmrelease canary-app
kubectl -n flux-system logs deployment/source-watcher --tail=100
kubectl -n flux-system logs deployment/helm-controller --tail=100
```

- Source not ready: check URL, branch and credentials. The default source uses public HTTPS.
- Generator not ready: check the copied paths and the source's ignore rules. Both `charts/` and `minikube/` must be included.
- Release not ready: check the ExternalArtifact feature gate, generated chart and Helm events.
- Pods not ready: inspect their events for image pull failures or application errors.
  This lab still uses the existing `ranjithka` application images.
- HTTP routing problems: inspect the ingress release and ingress class. This exercise
  retains the existing ingress-nginx chart pin; it does not upgrade that separate dependency.

## Migrating an existing installation

The release names remain `canary-dev` and `prd-dev`, the chart name remains `dev`,
and the namespace remains `default`. Production retains its deployment toleration
using the current `postRenderers.kustomize.patches` format.

Switching from `spec.chart` to `spec.chartRef` performs a Helm upgrade and Flux
cleans up the old generated HelmChart. When applying to an existing HelmRelease,
confirm that `spec.chart` has been removed; the two fields cannot coexist. If another
manager created it, remove the old field in that manager's configuration before applying.

This lab applies Flux resources with kubectl. Git supplies chart and values content;
changes to the Flux resource YAML itself require another `kubectl apply -k` unless
you separately configure a Flux Kustomization to manage these resources from Git.
Notifications and image automation remain separate exercises.
