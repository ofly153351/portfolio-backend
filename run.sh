#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
AI_DIR="$ROOT_DIR/ai-service"
AI_VENV="$AI_DIR/.venv"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "[error] missing command: $1" >&2
    exit 1
  fi
}

require_cmd go
require_cmd python3
require_cmd curl

if [[ -f "$ROOT_DIR/.env" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$ROOT_DIR/.env"
  set +a
fi

AI_PORT="${AI_SERVICE_PORT:-8000}"
GO_PORT="${PORT:-8080}"

if [[ ! -d "$AI_VENV" ]]; then
  echo "[setup] creating ai-service virtualenv"
  python3 -m venv "$AI_VENV"
fi

echo "[setup] installing ai-service dependencies"
"$AI_VENV/bin/pip" install -r "$AI_DIR/requirements.txt" >/dev/null

cleanup() {
  echo ""
  echo "[shutdown] stopping services..."
  if [[ -n "${GO_PID:-}" ]] && kill -0 "$GO_PID" 2>/dev/null; then
    kill "$GO_PID" 2>/dev/null || true
  fi
  if [[ -n "${AI_PID:-}" ]] && kill -0 "$AI_PID" 2>/dev/null; then
    kill "$AI_PID" 2>/dev/null || true
  fi
  wait || true
}
trap cleanup EXIT INT TERM

echo "[start] ai-service on :$AI_PORT"
(
  cd "$AI_DIR"
  "$AI_VENV/bin/uvicorn" app.main:app --host 0.0.0.0 --port "$AI_PORT" --reload
) &
AI_PID=$!

for _ in {1..60}; do
  if curl -fsS "http://127.0.0.1:$AI_PORT/health" >/dev/null 2>&1; then
    break
  fi
  sleep 0.5
done

if ! curl -fsS "http://127.0.0.1:$AI_PORT/health" >/dev/null 2>&1; then
  echo "[error] ai-service failed to become healthy on port $AI_PORT" >&2
  exit 1
fi

echo "[start] go-api on :$GO_PORT"
(
  cd "$ROOT_DIR"
  go run ./cmd/api
) &
GO_PID=$!

echo "[ready] services are running"
echo "        - Go API:     http://127.0.0.1:$GO_PORT"
echo "        - AI Service: http://127.0.0.1:$AI_PORT"
echo "[hint] press Ctrl+C to stop both"

# Portable wait loop for macOS bash 3.2 (no `wait -n`)
while true; do
  if ! kill -0 "$AI_PID" 2>/dev/null; then
    echo "[exit] ai-service stopped"
    break
  fi
  if ! kill -0 "$GO_PID" 2>/dev/null; then
    echo "[exit] go-api stopped"
    break
  fi
  sleep 1
done
