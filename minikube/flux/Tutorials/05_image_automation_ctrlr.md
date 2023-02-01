# Image Automation Controllers

The image-reflector-controller and image-automation-controller work together to update a Git repository when new container images are available.

-> The image-reflector-controller scans image repositories and reflects the image metadata in Kubernetes resources.

-> The image-automation-controller updates YAML files based on the latest images scanned, and commits the changes to a given Git repository.

<https://fluxcd.io/flux/guides/image-update/>


```bash
flux install --components-extra=image-reflector-controller,image-automation-controller
```

## Image Policy

<https://fluxcd.io/flux/components/image/imagepolicies/>

Image automation

Image Automation controller isn’t installed by default (since it’s an alpha feature) et you need to reconfigure Flux. The flux bootstrap command is idempotent and you can run it again with the new paramweters:

`--components-extra=image-reflector-controller,image-automation-controller`


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

`kubectl -n flux-system get secret flux-system -o json | jq '.data."identity.pub"' -r | base64 -d`

```bash
flux create secret git -n default git-secrets \
    --url=ssh://git@github.com/ranjith-ka/Devops \
    --private-key-file=/Users/ranjith.a/.ssh/id_rsa
```

```bash
$ kubectl get deploy canary-dev -o json | jq -r '.spec.template.spec.containers[0].image'
```
## Blogs to refer

<https://particule.io/en/blog/flux-auto-image-update/>