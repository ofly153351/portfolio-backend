#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="/home/ubuntu/portfolio/portfolio-backend"
AI_DIR="$PROJECT_DIR/ai-service"

GO_NAME="${GO_NAME:-go-api}"
AI_NAME="${AI_NAME:-ai-service}"

cd "$PROJECT_DIR"

pm2 delete "$GO_NAME" >/dev/null 2>&1 || true
pm2 delete "$AI_NAME" >/dev/null 2>&1 || true

pm2 start "go run cmd/api/main.go" --name "$GO_NAME" --cwd "$PROJECT_DIR"
pm2 start "uv run uvicorn app.main:app --host 0.0.0.0 --port 8000" --name "$AI_NAME" --cwd "$AI_DIR"

pm2 save
pm2 status
