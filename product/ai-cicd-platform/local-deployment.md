# Local Deployment

## Purpose

Run the AI-native CI/CD platform locally with the web app, inference service, PostgreSQL, and Redis.

## Prerequisites

- Docker Desktop or compatible container runtime
- Go 1.24+
- Node.js 20+

## Start Locally

```bash
docker compose -f docker-compose.local.yml up --build
```

## URLs

- Landing page: http://localhost:3000
- Dashboard: http://localhost:3000/dashboard
- Inference health: http://localhost:8080/healthz

## Runtime Notes

- The web app points to the inference service through the compose network.
- The landing page is the product entrypoint and the dashboard is the operational view.
- PostgreSQL and Redis are included for future contract-backed features.
