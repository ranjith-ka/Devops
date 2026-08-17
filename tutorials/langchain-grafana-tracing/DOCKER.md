# Docker Setup & Deployment Guide

## Overview

This project includes a complete Docker setup for running the LangChain tracing system with:
- Web UI for testing (Flask on port 5000)
- OpenTelemetry Collector for trace ingestion
- Grafana Tempo for trace storage
- Grafana for visualization

## Quick Start

### Prerequisites
- Docker Engine 20.10+
- Docker Compose v2.0+
- 4GB available disk space
- 2GB available RAM

### Build and Run

```bash
# Build Docker image (only needed first time)
docker compose build

# Start all services
docker compose up -d

# Check service status
docker compose ps
```

### Access the Application

| Service | URL | Credentials |
|---------|-----|-------------|
| **Web UI** | http://localhost:5000 | N/A |
| **Grafana** | http://localhost:3000 | admin / admin |
| **Tempo API** | http://localhost:3200 | N/A |
| **Loki API** | http://localhost:3100 | N/A |
| **OTLP Collector** | http://localhost:4317 (gRPC) | N/A |

## Web UI Features

### Tab 1: Question Answering
1. Navigate to http://localhost:5000
2. Enter a question in the textarea
3. Click "Send Question"
4. View the answer and trace ID

Example questions:
- "What problem does LangChain solve?"
- "How does tracing help debugging?"
- "Explain OpenTelemetry in simple terms"

### Tab 2: Compare Traces
1. Click the "Compare Traces" tab
2. Enter two trace IDs (e.g., `trace-a` and `trace-b`)
3. Click "Compare Traces"
4. View comparison report with:
   - Root cause identification
   - Span-by-span latency analysis
   - Processing flow
   - Recommendations

Use two 32-character trace IDs returned by successful requests in the
Question Answering tab. Comparisons load the recorded spans directly from Tempo.

## Monitoring with Grafana

### View Traces
1. Open http://localhost:3000
2. Login with admin/admin
3. Go to **Explore** tab
4. Select **Tempo** data source
5. Paste trace ID from Web UI
6. View full trace hierarchy

The comparison tab also queries Loki by each trace ID and displays correlated
node lifecycle, documentation retrieval, model, and memory events side by side.

### Configure Grafana (Optional)
- Default password: admin
- Change password on first login
- Add Tempo data source automatically configured

## Stopping Services

```bash
# Stop all containers
docker compose down

# Stop and remove volumes (clean slate)
docker compose down -v

# View logs
docker compose logs -f langchain-tracing-ui
docker compose logs -f otel-collector
```

## Environment Variables

Edit `.env` or modify `docker-compose.yaml`:

```yaml
environment:
  - SERVICE_NAME=langchain-tracing
  - SERVICE_VERSION=1.0.0
  - ENVIRONMENT=docker
  - OTLP_ENDPOINT=http://otel-collector:4317
```

## Troubleshooting

### Web UI not responding
```bash
# Check if container is running
docker compose ps langchain-tracing-ui

# View logs
docker compose logs langchain-tracing-ui

# Restart service
docker compose restart langchain-tracing-ui
```

### No traces appearing in Grafana
1. Ensure OTel Collector is running: `docker compose logs otel-collector`
2. Check Tempo is receiving traces: `docker compose logs tempo`
3. Verify OTLP_ENDPOINT is correct: `http://otel-collector:4317`

### High memory usage
- Reduce Tempo retention: Edit Tempo config in docker-compose.yaml
- Monitor with: `docker stats`

## Development Workflow

### Build Custom Image
```bash
# Rebuild after code changes
docker compose build --no-cache

# Run with updated code
docker compose up -d --force-recreate
```

### Local Testing (Without Docker)
```bash
# Install dependencies
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# Run CLI
python app.py "Your question"
python app.py --compare-trace-a trace-a --compare-trace-b trace-b

# Run Web UI (requires local OTLP endpoint)
python ui.py
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Docker Network (tracing)              │
├──────────────────┬──────────────────┬──────────────────┤
│                  │                  │                  │
│  Web UI (5000)   │  OTel Collector  │  Tempo (3200)    │
│  Flask App       │  (4317/4318)     │  Trace Storage   │
│                  │  Ingestion       │                  │
│                  │                  │                  │
└──────────────────┴──────────────────┴──────────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Grafana    │
                    │  (3000)      │
                    │ Visualization│
                    └──────────────┘
```

## Production Deployment

For production, consider:

1. **Security**
   - Use strong Grafana passwords
   - Implement authentication for Web UI
   - Use secrets for environment variables

2. **Performance**
   - Increase Tempo retention based on volume
   - Configure batch processor properly
   - Use dedicated storage for Tempo data

3. **Monitoring**
   - Set up Grafana alerts
   - Monitor container resource usage
   - Configure log aggregation

4. **Scaling**
   - Run multiple UI replicas behind load balancer
   - Scale OTel Collector independently
   - Use persistent storage for Tempo

See `ARCHITECTURE.md` for system design details.

## Support

For issues:
1. Check logs: `docker compose logs <service>`
2. Verify connectivity: `docker compose exec otel-collector curl http://tempo:3200/status`
3. Review configuration files
4. Check disk and memory: `docker system df`
