# Image Automation Controllers

The image-reflector-controller and image-automation-controller work together to update a Git repository when new container images are available.

-> The image-reflector-controller scans image repositories and reflects the image metadata in Kubernetes resources.
-> The image-automation-controller updates YAML files based on the latest images scanned, and commits the changes to a given Git repository.

## Image Policy

<https://fluxcd.io/flux/components/image/imagepolicies/>


```yaml
kind: ImagePolicy
spec:
  filterTags:
    pattern: '^RELEASE\.(?P<timestamp>.*)Z$'
    extract: '$timestamp'
  policy:
    alphabetical:
      order: asc
```

```yaml
kind: ImagePolicy
spec:
  policy:
    semver:
      range: '>=1.0.0 <2.0.0'
```

## Blogs to refer

<https://particule.io/en/blog/flux-auto-image-update/>