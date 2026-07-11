# Repository Structure

## Target Monorepo Layout

```text
repo/
├── apps/
│   └── web/
├── services/
│   ├── inference/
│   ├── api/
│   └── worker/
├── packages/
│   ├── contracts/
│   ├── ui/
│   └── config/
├── proto/
├── product/
│   └── ai-cicd-platform/
├── specs/
├── infra/
│   ├── terraform/
│   └── helm/
└── docs/
```

## Placement Rules

- Product strategy and UX live in `product/ai-cicd-platform/`.
- API and data contracts live in `specs/` and `proto/`.
- Reusable client types live in `packages/contracts`.
- Runtime services live in `services/`.
- Infrastructure lives in `infra/`.

## Current Status

- `apps/web` contains the Next.js UI shell.
- `services/inference` contains the Go inference API.
- `product/ai-cicd-platform` contains the product blueprint.
- `proto` already carries sample protobuf learning material.
