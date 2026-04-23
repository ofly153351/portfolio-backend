# AGENT CONTEXT: Portfolio CMS + AI Backend

## Structure Rules
Use this structure as the source of truth:

- `cmd/`: entrypoints
- `internal/app/`: HTTP bootstrap and route wiring
- `internal/config/`: environment-backed config
- `internal/database/`: MongoDB connection setup + pgvector connectivity for retrieval
- `internal/middleware/`: JWT auth and request guards
- `internal/modules/*`: bounded features for Portfolio CMS + AI only
- `internal/platform/httpx/`: shared JSON response helpers

## Module Reference
Each module should follow this baseline structure:

```text
./<module>
  errors.go
  handler.go
  model.go
  repository.go
  service.go
  util.go
```

Use this as the default pattern for new modules.

## Active Modules
- `auth`: authentication
- `content`: CMS content management
- `upload`: media upload flow
- `chat`: AI chat with RAG context
- `health`: runtime health checks

## Project Notes
- Primary DB: MongoDB
- Vector DB: pgvector (PostgreSQL)
- Framework: Go Fiber
- CORS must be configured from environment variables
- POS modules are intentionally removed from this repository

## API Module Map
- `auth` -> `/api/auth/*`
- `health` -> `/api/health`
- `content` -> `/api/admin/content`
- `upload` -> `/api/upload`
- `chat` -> `/api/chat`

## Runtime Env
- `.env.example` defines required variables
- `.env` is loaded at startup
- CORS env keys:
  - `CORS_ALLOW_ORIGINS`
  - `CORS_ALLOW_METHODS`
  - `CORS_ALLOW_HEADERS`
- Vector DB env keys:
  - `VECTOR_DB_HOST`
  - `VECTOR_DB_PORT`
  - `VECTOR_DB_USER`
  - `VECTOR_DB_PASSWORD`
  - `VECTOR_DB_NAME`
  - `VECTOR_DB_SSLMODE`
  - `VECTOR_DB_URL`

## Response Rules
- No reasoning. Just answer.
- Return only necessary output
- No explanation unless asked
- Always respond with minimal code
- No explanation unless asked
- Prefer short answers
