# Admin Authentication API

## POST /api/admin/login
Request:
```json
{
  "username": "admin",
  "password": "******"
}
```

Response 200:
```json
{
  "ok": true,
  "access_token": "token",
  "token_type": "Bearer",
  "user": {
    "id": "u_001",
    "username": "admin",
    "role": "admin"
  }
}
```

Response 401:
```json
{ "error": "invalid_credentials" }
```

Important:
- Backend returns token in response header: `Authorization: Bearer <token>`
- Backend also returns token in body (`access_token`, `token_type`)
- Frontend must send this header on all `/api/admin/*` endpoints

## POST /api/admin/logout
Response 200:
```json
{ "ok": true }
```

## GET /api/admin/me
Response 200:
```json
{
  "authenticated": true,
  "user": {
    "id": "u_001",
    "username": "admin",
    "role": "admin"
  }
}
```

Response 401:
```json
{ "authenticated": false }
```
