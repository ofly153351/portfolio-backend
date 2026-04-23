# GET /api/health

## Purpose
ตรวจสถานะ backend สำหรับ app startup/check readiness

## Current (Backend Now)
### Request
- Method: `GET`
- URL: `/api/health`

### Response
- Status: `200 OK`
```json
{
  "status": "ok",
  "app": "portfolio-backend"
}
```

## Planned Contract (Frontend Target)
คง format เดิมได้เลย (stable endpoint)

## Frontend Notes
- เรียกก่อนหน้า login/page render เพื่อเช็ค backend up/down
- ถ้า fail ให้ขึ้น maintenance state หรือ retry flow
