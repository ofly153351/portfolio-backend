# Backoffice Content API

รองรับ locale: `en`, `th`

## GET /api/admin/content?locale=en
Response 200:
```json
{
  "locale": "en",
  "version": 12,
  "updated_at": "2026-04-09T14:00:00Z",
  "content": {
    "technical": [
      {
        "id": "tech_1",
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
      "location": ""
    }
  }
}
```

## Technical (Separated Endpoints)

### GET /api/admin/technical?locale=en
Response 200:
```json
{
  "locale": "en",
  "version": 12,
  "updated_at": "2026-04-09T14:00:00Z",
  "items": [
    {
      "id": "tech_1",
      "title": "Go",
      "description": "Backend API",
      "icon": "http://localhost:9000/portfolio/technical/go.svg"
    }
  ]
}
```

### POST /api/admin/technical?locale=en
Request:
```json
{
  "title": "Redis",
  "description": "Cache and memory store",
  "icon": "http://localhost:9000/portfolio/technical/redis.svg"
}
```

### PUT /api/admin/technical/{id}?locale=en
Request:
```json
{
  "title": "Redis",
  "description": "Cache, queue, memory",
  "icon": "http://localhost:9000/portfolio/technical/redis.svg"
}
```

### DELETE /api/admin/technical/{id}?locale=en
Response 200:
```json
{
  "ok": true,
  "version": 13,
  "updated_at": "2026-04-09T14:10:00Z",
  "deleted_id": "tech_1"
}
```

## PUT /api/admin/content?locale=en
Request:
```json
{
  "version": 12,
  "content": {
    "technical": [],
    "projects": [
      {
        "id": "proj_1",
        "tag": "AI",
        "title": "Portfolio CMS",
        "description": "Backoffice + AI chat",
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
      "location": ""
    }
  }
}
```

Response 200:
```json
{
  "ok": true,
  "version": 13,
  "updated_at": "2026-04-09T14:10:00Z"
}
```

Response 409:
```json
{
  "error": "version_conflict",
  "current_version": 13
}
```

## POST /api/admin/content/publish?locale=en
Response 200:
```json
{
  "ok": true,
  "published_version": 13,
  "published_at": "2026-04-09T14:15:00Z"
}
```

## GET /api/admin/content/history?locale=en
Response 200:
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

## GET /api/content?locale=en
Public endpoint for frontend portfolio.

Response 200:
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

## Validation Rules
- `locale` only `en` or `th`
- `technical[].title` and `projects[].title` are required
- `technical[].icon` must be valid `http/https` URL when provided
- `projects[].image` must be valid `http/https` URL when provided
- `projects[].images[]` must be valid `http/https` URL when provided
- `about` max 5000 chars
- `technical[].description` max 2000 chars
- `projects[].description` max 3000 chars

## Persistence
- Draft/Published state is split into MongoDB collections:
  - `portfolio_projects`
  - `portfolio_technical`
  - `portfolio_info`
- Edit history is stored in MongoDB collection `portfolio_content_history`
- Legacy `portfolio_content` is migrated automatically on startup
