# Media Upload API

## POST /api/admin/upload
Upload media and return URL(s) for `projects[].image`/`projects[].images`.

Request:
- Method: `POST`
- Content-Type: `multipart/form-data`
- Single file field: `file`
- Multiple files field: `files` (repeat key)

Response 200:
```json
{
  "ok": true,
  "url": "http://localhost:9000/portfolio/projects/20260409_<uuid>.jpg",
  "urls": [
    "http://localhost:9000/portfolio/projects/20260409_<uuid>.jpg",
    "http://localhost:9000/portfolio/projects/20260409_<uuid>.gif"
  ]
}
```

Response 400:
```json
{ "error": "file_required" }
```

Response 400:
```json
{ "error": "invalid_file_type" }
```

Response 401:
```json
{ "authenticated": false }
```

Notes:
- Endpoint requires `Authorization: Bearer <token>`
- Backend uploads binary file to MinIO bucket (`MINIO_BUCKET`) and returns object URL
- Allowed types: `jpg`, `jpeg`, `png`, `webp`, `gif`, `svg`

## POST /api/admin/technical/upload
Upload icon/media for `technical[].icon`.

Request:
- Method: `POST`
- Content-Type: `multipart/form-data`
- Single file field: `file`
- Multiple files field: `files` (repeat key)

Response 200:
```json
{
  "ok": true,
  "url": "http://localhost:9000/portfolio/technical/20260409_<uuid>.svg",
  "urls": [
    "http://localhost:9000/portfolio/technical/20260409_<uuid>.svg"
  ]
}
```
