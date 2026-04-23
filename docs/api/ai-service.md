# AI Service API (LangChain)

Base URL: `http://localhost:8000`

## Scope (Streaming Policy)
- AI service exposes **non-streaming** endpoints only:
  - `POST /chat`
  - `POST /embed`
- Streaming for client is handled only by Go WebSocket endpoint:
  - `ws://localhost:8080/api/chat/ws`

## Quick Run
```bash
cd ai-service
source .venv/bin/activate
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

## Provider Config

### OpenAI
```env
API_SERVICE=openai
OPENAI_API_KEY=sk-...
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_CHAT_MODEL=gpt-4o-mini
OPENAI_EMBED_MODEL=text-embedding-3-small
OPENAI_TIMEOUT_SECONDS=60
```

### Ollama
```env
API_SERVICE=ollama
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL=gemma3:1b
OLLAMA_EMBED_MODEL=nomic-embed-text
OLLAMA_TIMEOUT_SECONDS=300
```

## GET /health
Response:
```json
{ "status": "ok", "service": "ai-service" }
```

## GET /admin/chat-memory/{session_id}
ใช้สำหรับ debug memory ของ session

Postman:
- Method: `GET`
- URL: `http://localhost:8000/admin/chat-memory/sess_001`

Response:
```json
{
  "session_id": "sess_001",
  "memory_store": "redis",
  "count": 2,
  "turns": [
    { "role": "user", "content": "hello" },
    { "role": "assistant", "content": "hi" }
  ]
}
```

## DELETE /admin/chat-memory/{session_id}
ใช้สำหรับลบ memory ของ session

Postman:
- Method: `DELETE`
- URL: `http://localhost:8000/admin/chat-memory/sess_001`

Response:
```json
{
  "session_id": "sess_001",
  "memory_store": "redis",
  "deleted": true,
  "cleared_count": 2
}
```

## POST /chat
Postman:
- Method: `POST`
- URL: `http://localhost:8000/chat`
- Header: `Content-Type: application/json`

Request:
```json
{
  "message": "hello",
  "session_id": "sess_001",
  "top_k": 5,
  "lang": "th"
}
```

Response:
```json
{
  "answer": "...model output...",
  "sources": [],
  "session_id": "sess_001",
  "provider": "openai",
  "usage": {
    "prompt_tokens": 120,
    "completion_tokens": 80,
    "total_tokens": 200
  }
}
```

Notes:
- `provider` tells which backend served the request (`openai` or `ollama`).
- `usage` is per-turn token usage from provider response.
- `AI_SYSTEM_PROMPT` is injected on every chat call from root `.env`.
- `session_id` เดิม = ใช้ memory เดิม, `session_id` ใหม่ = เริ่มบทสนทนาใหม่
- จำนวนรอบ memory คุมด้วย `MEMORY_MAX_TURNS` (default `12`)
- memory backend ตอนนี้ใช้ Redis (`MEMORY_STORE=redis`)

## POST /embed
Postman:
- Method: `POST`
- URL: `http://localhost:8000/embed`
- Header: `Content-Type: application/json`

Request:
```json
{
  "id": "cnt_001",
  "title": "POS System",
  "content": "Full-stack retail POS",
  "type": "project"
}
```

Response:
```json
{
  "id": "cnt_001",
  "dimensions": 1536,
  "status": "embedded"
}
```
