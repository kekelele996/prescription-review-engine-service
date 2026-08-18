#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="${RUNTIME_ENV_FILE:-.env.example}"

cleanup() {
  docker compose --env-file "$ENV_FILE" down --volumes --remove-orphans >/dev/null 2>&1 || true
}

trap cleanup EXIT INT TERM
cleanup
docker compose --env-file "$ENV_FILE" up --build backend
