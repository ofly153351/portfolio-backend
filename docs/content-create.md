# Content Create API

รองรับ `locale`: `en`, `th`

## Data Shape

`content` มี 3 ส่วน:
- `technical[]`
- `projects[]`
- `portfolioInfo`

ตัวอย่าง `projects[]` (อัปเดตล่าสุด):
```json
{
  "id": "proj_1",
  "index": 0,
  "tag": "AI",
  "title": "Portfolio CMS",
  "description": "Backoffice + AI chat",
  "repoUrl": "https://github.com/username/portfolio-cms",
  "projectUrl": "https://portfolio.example.com",
  "image": "http://localhost:9000/portfolio/projects/cover.jpg",
  "images": [
    "http://localhost:9000/portfolio/projects/cover.jpg",
    "http://localhost:9000/portfolio/projects/demo.gif"
  ]
}
```

## 1) GET `/api/admin/content?locale=en`
Response `200`:
```json
{
  "locale": "en",
  "version": 12,
  "updated_at": "2026-04-09T14:00:00Z",
  "content": {
    "technical": [
      {
        "id": "tech_1",
        "index": 0,
        "title": "Go",
        "description": "Backend API",
        "icon": "http://localhost:9000/portfolio/technical/go.svg"
      }
    ],
    "projects": [],
    "portfolioInfo": {
      "ownerName": "",
      "title": "",
      "subtitle": "",
      "about": "",
      "contactEmail": "",
      "contactPhone": "",
      "location": "",
      "github": "",
      "linkedin": "",
      "instagram": ""
    }
  }
}
```

## 2) PUT `/api/admin/content?locale=en`
Request:
```json
{
  "version": 12,
  "content": {
    "technical": [],
    "projects": [
      {
        "id": "proj_1",
        "index": 0,
        "tag": "AI",
        "title": "Portfolio CMS",
        "description": "Backoffice + AI chat",
        "repoUrl": "https://github.com/username/portfolio-cms",
        "projectUrl": "https://portfolio.example.com",
        "image": "http://localhost:9000/portfolio/projects/cover.jpg",
        "images": [
          "http://localhost:9000/portfolio/projects/cover.jpg",
          "http://localhost:9000/portfolio/projects/demo.gif"
        ]
      }
    ],
    "portfolioInfo": {
      "ownerName": "",
      "title": "",
      "subtitle": "",
      "about": "",
      "contactEmail": "",
      "contactPhone": "",
      "location": "",
      "github": "https://github.com/username",
      "linkedin": "https://www.linkedin.com/in/username",
      "instagram": "https://www.instagram.com/username"
    }
  }
}
```

Response `200`:
```json
{
  "ok": true,
  "version": 13,
  "updated_at": "2026-04-09T14:10:00Z"
}
```

Response `409`:
```json
{
  "error": "version_conflict",
  "current_version": 13
}
```

## 3) POST `/api/admin/content/publish?locale=en`
ใช้สำหรับนำ `draft` ล่าสุดของ locale นั้นไปทับเป็น `published` เพื่อให้ Public endpoint (`/api/content`) อ่านข้อมูลล่าสุดได้

Headers:
```http
Authorization: Bearer <access_token> 
Content-Type: application/json
```

Request Body: ไม่มี

Example:
```bash
curl -X POST "http://localhost:8080/api/admin/content/publish?locale=en" \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json"
```

ตัวอย่างภาษาไทย:
```bash
curl -X POST "http://localhost:8080/api/admin/content/publish?locale=th" \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json"
```

Response `200`:
```json
{
  "ok": true,
  "published_version": 13,
  "published_at": "2026-04-09T14:15:00Z"
}
```

Response `400`:
```json
{
  "error": "invalid_locale"
}
```

Response `401`:
```json
{
  "authenticated": false
}
```

Response `500`:
```json
{
  "error": "internal_error"
}
```

## 4) GET `/api/admin/content/history?locale=en`
Response `200`:
```json
{
  "locale": "en",
  "history": [
    {
      "locale": "en",
      "version": 13,
      "updated_by": "admin",
      "updated_at": "2026-04-09T14:10:00Z",
      "content": {
        "technical": [],
        "projects": [],
        "portfolioInfo": {}
      }
    }
  ]
}
```

## 5) GET `/api/content?locale=en`
Public endpoint สำหรับ frontend portfolio

Headers:
```http
X-Public-Token: <public_token>
```

โดย token ได้จาก `GET /api/public/token`

Response `200`:
```json
{
  "locale": "en",
  "content": {
    "technical": [],
    "projects": [],
    "portfolioInfo": {}
  }
}
```

Response `401`:
```json
{
  "error": "public_token_invalid"
}
```

## Technical Endpoints

### GET `/api/admin/technical?locale=en`
Response `200`:
```json
{
  "locale": "en",
  "version": 12,
  "updated_at": "2026-04-09T14:00:00Z",
  "items": [
    {
      "id": "tech_1",
      "index": 0,
      "title": "Go",
      "description": "Backend API",
      "icon": "http://localhost:9000/portfolio/technical/go.svg"
    }
  ]
}
```

### POST `/api/admin/technical?locale=en`
Request:
```json
{
  "index": 0,
  "title": "Redis",
  "description": "Cache and memory store",
  "icon": "http://localhost:9000/portfolio/technical/redis.svg"
}
```

### PUT `/api/admin/technical/{id}?locale=en`
Request:
```json
{
  "index": 1,
  "title": "Redis",
  "description": "Cache, queue, memory",
  "icon": "http://localhost:9000/portfolio/technical/redis.svg"
}
```

### DELETE `/api/admin/technical/{id}?locale=en`
Response `200`:
```json
{
  "ok": true,
  "version": 13,
  "updated_at": "2026-04-09T14:10:00Z",
  "deleted_id": "tech_1"
}
```

## Validation Rules

- `locale` ต้องเป็น `en` หรือ `th`
- `technical[].title` และ `projects[].title` จำเป็นต้องมี
- `technical[].index` ต้องเป็นจำนวนเต็ม >= 0
- ระบบจะ normalize `technical[].index` ใหม่ตามลำดับใน array ก่อนบันทึก (0..n-1)
- `projects[].index` ต้องเป็นจำนวนเต็ม >= 0
- ระบบจะ normalize `projects[].index` ใหม่ตามลำดับใน array ก่อนบันทึก (0..n-1)
- `projects[].repoUrl` ต้องเป็น URL `http/https` (ถ้าส่งมา)
- `technical[].icon` ต้องเป็น URL `http/https` (ถ้าส่งมา)
- `projects[].projectUrl` ต้องเป็น URL `http/https` (ถ้าส่งมา)
- `projects[].image` ต้องเป็น URL `http/https` (ถ้าส่งมา)
- `projects[].images[]` ต้องเป็น URL `http/https` (ถ้าส่งมา)
- `portfolioInfo.about` ยาวไม่เกิน 5000 ตัวอักษร
- `portfolioInfo.github` ต้องเป็น URL `http/https` (ถ้าส่งมา)
- `portfolioInfo.linkedin` ต้องเป็น URL `http/https` (ถ้าส่งมา)
- `portfolioInfo.instagram` ต้องเป็น URL `http/https` (ถ้าส่งมา)
- `technical[].description` ยาวไม่เกิน 2000 ตัวอักษร
- `projects[].description` ยาวไม่เกิน 3000 ตัวอักษร

## Persistence

- Draft/Published แยกเก็บใน MongoDB collections:
  - `portfolio_projects`
  - `portfolio_technical`
  - `portfolio_info`
- History เก็บใน `portfolio_content_history`
