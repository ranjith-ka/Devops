# Tutorial Agent (Go)

Go-based AI tutorial assistant for:
- building learning paths
- creating step-by-step lessons
- generating hands-on exercises
- saving each generated tutorial as a `.md` file under `minikube/tutorial-ai/`

## 1) Configure environment

```bash
cd agents/tutorial_agent
cp .env.example .env
```

Required:
- `GITHUB_MODELS_TOKEN`

Optional:
- `GITHUB_MODELS_MODEL` (default: `openai/gpt-5-mini`)
- `GITHUB_MODELS_ENDPOINT` (default: `https://models.github.ai/inference/chat/completions`)
- `PORT` (default: `8088`)
- `TUTORIAL_OUTPUT_DIR` (optional override for markdown output folder)

## 2) Run

From repo root:

```bash
go run ./agents/tutorial_agent
```

## 3) API

### Health

```bash
curl http://localhost:8088/health
```

### Tutorial

```bash
curl -X POST http://localhost:8088/tutorial \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "lesson",
    "topic": "kubernetes deployments",
    "level": "beginner",
    "outputFile": "k8s-deployments-lesson.md",
    "message": "Teach me with an example and a mini practice task"
  }'
```

The agent writes a markdown file in `minikube/tutorial-ai/` by default and returns:
- `filePath`
- `fileName`
- `savedIn`
- `timestamp`

`mode` options:
- `plan`
- `lesson`
- `exercise`
- default: general tutorial guidance

## 4) VS Code task

Run task: `Run Tutorial Agent HTTP Server` from [agents/tutorial_agent/.vscode/tasks.json](./.vscode/tasks.json).
