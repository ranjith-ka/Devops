# Kustomize

<https://www.densify.com/kubernetes-tools/kustomize>

Check the documentation in the above link, hope it will work after sometime as well.

Overline for Kustomize

```bash
├── base
  │   ├── deployment.yaml
  │   ├── hpa.yaml
  │   ├── kustomization.yaml
  │   └── service.yaml
  └── overlays
      ├── dev
      │   ├── hpa.yaml
      │   └── kustomization.yaml
      ├── production
      │   ├── hpa.yaml
      │   ├── kustomization.yaml
      │   ├── rollout-replica.yaml
      │   └── service-loadbalancer.yaml
      └── staging
          ├── hpa.yaml
          ├── kustomization.yaml
          └── service-nodeport.yaml
```