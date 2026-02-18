# Code Architect Agent (Go)

Go-based AI partner for:
- Code editing guidance (diff-first responses)
- New architecture/design discussions
- Migration and refactor planning

## 1) Configure environment

Create a local env file:

```bash
cp .env.example .env
```

This agent is configured for GitHub Models only.

### GitHub Models

- `GITHUB_MODELS_TOKEN` (GitHub PAT)
- optional `GITHUB_MODELS_MODEL` (default: `openai/gpt-5-mini`)
- optional `GITHUB_MODELS_ENDPOINT` (default: `https://models.github.ai/inference/chat/completions`)

## 2) Run

From repo root:

```bash
go run ./agents/code_architect_agent
```

Server starts on `PORT` (default `8087`).

Note: `go run ./agents/code_architect_agent` now auto-loads env values from either:
- `.env` (repo root)
- `agents/code_architect_agent/.env`

## 3) API

### Health

```bash
curl http://localhost:8087/health
```

### Chat

```bash
curl -X POST http://localhost:8087/chat \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "architecture",
    "context": "Service currently monolithic with MongoDB",
    "message": "Propose target architecture for scaling and safer deployments"
  }'
```

`mode` can be:
- `edit`
- `architecture`
- `general` (default)

## 4) VS Code task

Run task: `Run Agent HTTP Server` from [agents/code_architect_agent/.vscode/tasks.json](./.vscode/tasks.json).
