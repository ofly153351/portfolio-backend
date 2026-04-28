#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="/home/ubuntu/portfolio/portfolio-backend"
AI_DIR="$PROJECT_DIR/ai-service"

GO_NAME="${GO_NAME:-go-api}"
AI_NAME="${AI_NAME:-ai-service}"

export PATH="$PATH:/usr/bin:/usr/local/bin"
if [ -s "$HOME/.nvm/nvm.sh" ]; then
  . "$HOME/.nvm/nvm.sh"
  nvm use --lts >/dev/null 2>&1 || true
fi

PM2_BIN="$(command -v pm2 || true)"
if [ -z "$PM2_BIN" ]; then
  echo "pm2 not found in PATH"
  exit 127
fi

cd "$PROJECT_DIR"

"$PM2_BIN" delete "$GO_NAME" >/dev/null 2>&1 || true
"$PM2_BIN" delete "$AI_NAME" >/dev/null 2>&1 || true

"$PM2_BIN" start "go run cmd/api/main.go" --name "$GO_NAME" --cwd "$PROJECT_DIR"
"$PM2_BIN" start "uv run uvicorn app.main:app --host 0.0.0.0 --port 8000" --name "$AI_NAME" --cwd "$AI_DIR"

"$PM2_BIN" save
"$PM2_BIN" status
