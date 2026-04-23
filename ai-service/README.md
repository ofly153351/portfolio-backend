# AI Service (LangChain)

Python service for AI chat and embedding APIs.

## Endpoints
- `GET /health`
- `POST /chat`
- `POST /embed`
- `GET /admin/chat-memory/{session_id}` (debug memory)
- `DELETE /admin/chat-memory/{session_id}` (clear memory)

## Streaming Policy
- AI service does not expose streaming endpoint.
- Client realtime streaming is provided by Go API WebSocket:
  - `ws://localhost:8080/api/chat/ws`

## Run (local)
```bash
cd ai-service
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

## Required env
- `API_SERVICE` (`openai` or `ollama`, default `openai`)
- `OPENAI_API_KEY`
- `OPENAI_BASE_URL` (OpenAI/OpenRouter)
- `OPENAI_CHAT_MODEL`
- `OPENAI_EMBED_MODEL`
- `OPENAI_TIMEOUT_SECONDS` (optional, default `60`)
- `OLLAMA_BASE_URL` (when `API_SERVICE=ollama`)
- `OLLAMA_MODEL` (when `API_SERVICE=ollama`)
- `OLLAMA_EMBED_MODEL` (when `API_SERVICE=ollama`)
- `OLLAMA_TIMEOUT_SECONDS` (optional, default `300`)
- `OPENROUTER_HTTP_REFERER` (optional)
- `OPENROUTER_X_TITLE` (optional)
- `AI_SYSTEM_PROMPT` (optional)
- `MEMORY_MAX_TURNS` (optional, default `12`)
- `MEMORY_STORE` (`redis` or `memory`, default `redis`)
- `REDIS_URL` (required when `MEMORY_STORE=redis`)
- `REDIS_KEY_PREFIX` (optional, default `portfolio:chat:`)
- `REDIS_TIMEOUT_SECONDS` (optional, default `5`)

## Production Structure
ดูรายละเอียดที่ไฟล์ `structure.md`
