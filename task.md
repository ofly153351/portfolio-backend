# Task: Portfolio Backend Split into 2 Services

## Goal
ทำระบบแบบ monorepo แต่แยก runtime ชัดเจนเป็น 2 service:
1. Go API (CMS + upload + gateway)
2. AI Service (Python + LangChain for chat/embed)

## Target Layout
```text
portfolio-backend/
├── cmd/
│   └── api/main.go                # Go API entrypoint
├── internal/                      # Go API code
├── ai-service/
│   ├── app/main.py                # FastAPI + LangChain entrypoint
│   ├── requirements.txt
│   └── README.md
├── docs/api/
├── docker-compose.yml             # shared infra only
├── .env
└── .env.example
```

## Service Responsibilities

### 1) Go API Service
- CMS endpoints (`/api/admin/content`)
- Upload endpoint (`/api/upload`)
- Health/auth endpoints
- Chat gateway endpoint (`/api/chat`) that forwards to AI service
- หลัง save content ให้เรียก AI service `/embed` (phase ถัดไป)

### 2) AI Service (LangChain)
- `POST /chat`:
  - รับ prompt
  - เรียก LLM ผ่าน OpenAI-compatible provider
  - คืน answer
- `POST /embed`:
  - รับ content
  - สร้าง embedding
  - คืน metadata สำหรับนำไปเขียน vector store

## Communication Contract
- Go API -> AI Service base URL จาก `AI_SERVICE_URL`
- Default: `http://localhost:8000`
- Go `/api/chat` จะ call AI `/chat`

## Shared Infrastructure
- MongoDB: content storage
- MinIO: image/object storage
- pgvector: retrieval/vector indexing

## Environment
- Go API ต้องใช้:
  - `AI_SERVICE_URL`
  - `MONGO_URI`, `MONGO_DB`
  - MinIO, CORS config
- AI Service ต้องใช้:
  - `OPENAI_API_KEY`
  - `OPENAI_BASE_URL`
  - `OPENAI_CHAT_MODEL`
  - `OPENAI_EMBED_MODEL`
  - `OPENROUTER_HTTP_REFERER` (optional)
  - `OPENROUTER_X_TITLE` (optional)

## RAG Flow (target)
1. Go CMS save content
2. Go call AI `/embed`
3. Store vectors in pgvector
4. Chat request -> AI retrieve context + generate answer

## Persona
```text
You are Peerapat.

Style:
- practical
- opinionated
- dev mindset

Context:
{retrieved_docs}
```
