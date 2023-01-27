# Kustomize and Helm controller

There are few controller CRD's get installed with flux installation, to manage the files we can use any of the controllers as per project needs.

## kustomize controllers

`source.toolkit.fluxcd.io/v1beta2`  - Source to create for Single source ot Truth
`kustomize.toolkit.fluxcd.io/v1beta2` - Kustomize contoller to pull the files and apply to the NS or cluster

<https://fluxcd.io/flux/components/kustomize/kustomization/>  - Detailed docs in there to manage kustomize the files
 

```yaml
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: GitRepository
metadata:
  name: podinfo
  namespace: default
spec:
  interval: 5m
  url: https://github.com/stefanprodan/podinfo
  ref:
    branch: master
---
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  name: podinfo
  namespace: default
spec:
  interval: 10m
  targetNamespace: default
  sourceRef:
    kind: GitRepository
    name: podinfo
  path: "./kustomize"
  prune: true
```

### Helm controllers

HelmRelease resources has a built-in Kustomize compatible Post Renderer, which provides the following Kustomize directives:
<https://fluxcd.io/flux/components/helm/helmreleases/#post-renderers>

```yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: metrics-server
  namespace: kube-system
spec:
  interval: 1m
  chart:
    spec:
      chart: metrics-server
      version: "5.3.4"
      sourceRef:
        kind: HelmRepository
        name: bitnami
        namespace: kube-system
      interval: 1m
  postRenderers:
    # Instruct helm-controller to use built-in "kustomize" post renderer.
    - kustomize:
        # Array of inline strategic merge patch definitions as YAML object.
        # Note, this is a YAML object and not a string, to avoid syntax
        # indention errors.
        patchesStrategicMerge:
          - kind: Deployment
            apiVersion: apps/v1
            metadata:
              name: metrics-server
            spec:
              template:
                spec:
                  tolerations:
                    - key: "workload-type"
                      operator: "Equal"
                      value: "cluster-services"
                      effect: "NoSchedule"
        # Array of inline JSON6902 patch definitions as YAML object.
        # Note, this is a YAML object and not a string, to avoid syntax
        # indention errors.
        patchesJson6902:
          - target:
              version: v1
              kind: Deployment
              name: metrics-server
            patch:
              - op: add
                path: /spec/template/priorityClassName
                value: system-cluster-critical
        images:
          - name: docker.io/bitnami/metrics-server
            newName: docker.io/bitnami/metrics-server
            newTag: 0.4.1-debian-10-r54
```