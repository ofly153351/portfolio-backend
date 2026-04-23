# API Docs (Frontend Prompt Ready)

## Backoffice CMS
- `auth-login.md` -> `/api/admin/login`, `/api/admin/logout`, `/api/admin/me`
- `content-create.md` -> `/api/admin/content`, `/api/admin/content/publish`, `/api/admin/content/history`, `/api/content`
- `content-create.md` -> `/api/admin/technical` (GET/POST/PUT/DELETE) แยกจาก content ก้อนใหญ่
- `upload-image.md` -> `/api/admin/upload`, `/api/admin/technical/upload`

## AI + Chat
- `api/chat-websocket.md` -> `ws://localhost:8080/api/chat/ws`
- `api/ai-service.md` -> `POST /chat`, `POST /embed`, memory debug endpoints

## Security Contract
- `/api/admin/*` requires `Authorization: Bearer <token>` except `/api/admin/login`
- Login token comes from `Authorization` response header of `/api/admin/login`
