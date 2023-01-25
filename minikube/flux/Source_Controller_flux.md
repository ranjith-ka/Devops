### Flux Again

1. ## Source Controller

    - Buckets
    - Git Repo
    - Helm charts
    - helm repo
    - OCI repo (skipped no plan as of now)
    
### Buckets:
 - API defines a Source to produce an Artifact for objects from storage solutions like Amazon S3, Google Cloud Storage buckets, or any other solution with a S3 compatible API such as Minio, Alibaba Cloud OSS and others.
<https://fluxcd.io/flux/components/source/buckets/> - Example API's

`flux reconcile source bucket <bucket-name>` 

- reconcile
- suspend Bucket
- Resume Bucket

#### GIT Repo:

The GitRepository API defines a Source to produce an Artifact for a Git repository revision.
<https://fluxcd.io/flux/components/source/gitrepositories/>

```yaml
---
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: GitRepository
metadata:
  name: podinfo
  namespace: default
spec:
  interval: 5m0s
  url: https://github.com/stefanprodan/podinfo
  ref:
    branch: master
```

- Interval (time to sync)
- Timeout (timeout for long process)
- Reference
- branch
- Tag
- Semver (example below)


```yaml
apiVersion: source.toolkit.fluxcd.io/v1alpha1
kind: GitRepository
metadata:
  name: podinfo
spec:
  interval: 1m
  url: https://github.com/stefanprodan/podinfo
  ref:
    semver: ">=3.1.0-rc.1 <3.2.0"
```

### Helm charts

The HelmChart API defines a Source to produce an Artifact for a Helm chart archive with a set of specific configurations.

```yaml
spec:
  chart:
    spec:
      chart: podinfo
      ...
      valuesFiles:
        - values.yaml
        - values-production.yaml
```
NOTE: If the reconcile strategy is ChartVersion and the source reference is a GitRepository or a Bucket, no new chart artifact is produced on updates to the source unless the version in Chart.yaml is incremented. To produce new chart artifact on change in source revision, set the reconcile strategy to `Revision`


### Helm Repo 

There are 2 Helm repository types defined by the HelmRepository API:

Helm HTTP/S repository, which defines a Source to produce an Artifact for a Helm repository index YAML (index.yaml).

OCI Helm repository, which defines a source that does not produce an Artifact. Instead a validation of the Helm repository is performed and the outcome is reported in the .status.conditions field.

