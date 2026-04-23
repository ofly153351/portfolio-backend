# Portfolio CMS + AI Backend

Go Fiber backend scaffolded with Clean Architecture principles for a Portfolio CMS with AI integration.

## Project Structure

```text
.
├── cmd/                        # entrypoints
├── internal/
│   ├── app/                    # HTTP bootstrap and route wiring
│   ├── config/                 # environment-backed configuration
│   ├── database/               # MongoDB connection (+ pgvector for retrieval)
│   ├── middleware/             # JWT auth and request guards
│   ├── modules/                # bounded features
│   │   ├── auth/
│   │   ├── content/
│   │   ├── upload/
│   │   ├── chat/
│   │   └── health/
│   └── platform/httpx/         # shared JSON response helpers
├── docker-compose.yml
└── README.md
```

## API Routes

- `POST /api/admin/login`
- `POST /api/admin/logout`
- `GET /api/admin/me`
- `GET /api/health`
- `GET /api/admin/content?locale=en|th`
- `PUT /api/admin/content?locale=en|th`
- `POST /api/admin/content/publish?locale=en|th`
- `GET /api/admin/content/history?locale=en|th`
- `GET /api/admin/technical?locale=en|th`
- `POST /api/admin/technical?locale=en|th`
- `PUT /api/admin/technical/:id?locale=en|th`
- `DELETE /api/admin/technical/:id?locale=en|th`
- `POST /api/admin/upload`
- `POST /api/admin/technical/upload`
- `GET /api/content?locale=en|th`
- `GET ws://localhost:8080/api/chat/ws`

## Infrastructure

- MongoDB: source content storage
- pgvector (PostgreSQL): vector similarity search for RAG
- MinIO: media/object storage

## Run

```bash
cp .env.example .env
docker compose up -d
```

Run Go API service:

```bash
go mod tidy
go build ./...
go run ./cmd/api
```

Or run both Go API + AI service together:
```bash
./run.sh
```

Run AI service (LangChain):

```bash
cd ai-service
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

Go API will call AI service via `AI_SERVICE_URL` (default `http://localhost:8000`).

For OpenRouter, set:
- `OPENAI_BASE_URL=https://openrouter.ai/api/v1`
- `OPENAI_CHAT_MODEL` to an OpenRouter model id

## CORS

CORS is configured from `.env`:
- `CORS_ALLOW_ORIGINS`
- `CORS_ALLOW_METHODS`
- `CORS_ALLOW_HEADERS`

## Upload Storage

- `POST /api/admin/upload` stores files in MinIO (`MINIO_BUCKET`)
- URL returned is based on `MINIO_PUBLIC_BASE_URL`
