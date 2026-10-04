# Web App

Next.js + Tailwind + shadcn-ready shell for the AI-native CI/CD Intelligence Platform.

## Flux application onboarding

From the repository root:

```bash
npm ci --prefix apps/web
make flux-ui
```

Open http://localhost:3100/platform-onboarding. Enter an application name,
environment, existing image repository and versioned tag, replicas, container port,
and optional hostname. Generate the configuration to review each file, then download
the ZIP. All generation runs locally; no GitHub credentials are needed.

The bundle contains environment values, a Namespace, ArtifactGenerator,
HelmRelease, Kustomize entry point and `ONBOARDING.md`. Extract it at the repository
root and add the resource entry from `ONBOARDING.md` to
`minikube/flux/staging/kustomization.yaml`, then open a PR. After merge, the existing
`devops-staging` Flux Kustomization manages the application. The platform must already
be initialized with `make flux` and `make flux-app`.

Each application/environment has its own namespace and generated artifact.
Upgrade by editing the generated values file's `image.tag` and merging a new PR.
The UI creates downloadable files; it does not create PRs or apply to Kubernetes.
It uses `charts/dev`, whose container command is `serve`; images must support that
command. Configure imagePullSecrets separately for private registries.

Validation:

```bash
npm --prefix apps/web run test:onboarding
npm --prefix apps/web run build
```
