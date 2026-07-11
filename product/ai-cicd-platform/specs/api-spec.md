# API Specification

## Inference APIs

All endpoints return a consistent envelope with request id, model metadata, evidence, recommendations, and structured artifacts.

### Common Response

```json
{
  "request_id": "uuid",
  "model": "gpt-5.4-mini",
  "prompt_template": "delivery-artifact-v1",
  "summary": "string",
  "confidence": 0.92,
  "evidence": ["string"],
  "recommendations": ["string"],
  "artifacts": ["string"]
}
```

### Endpoints

- `POST /chat`
- `POST /analyze-pipeline`
- `POST /deployment-summary`
- `POST /root-cause`
- `POST /optimize`
- `POST /security-review`
- `POST /explain-workflow`
- `POST /generate-pipeline`
- `POST /generate-terraform`
- `POST /generate-kubernetes`

## Request Envelope

```json
{
  "org_id": "string",
  "repository": "string",
  "pipeline": "string",
  "deployment": "string",
  "prompt": "string",
  "context": {
    "branch": "string",
    "environment": "string",
    "time_range": "string"
  }
}
```

## Behavior Rules

- Every response must cite evidence if the request references logs, PRs, deployments, or metrics.
- Artifact-generation endpoints must validate structure before returning output.
- `security-review` must rank findings by severity and confidence.
- `root-cause` must return ranked hypotheses rather than a single guess.
- `optimize` must split recommendations into quick wins and structural changes.
- All routes must be tenant-scoped and policy-aware.

## Contract Notes

- Use REST for UI and external integration compatibility.
- Keep a stable JSON envelope for web, CLI, and automation clients.
- Version prompts independently from API schema.
