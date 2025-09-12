#!/usr/bin/env bash
set -e
export COMPOSE_HTTP_TIMEOUT=200
docker compose -f docker-compose.demo.yml up --build
