# AI Service Production Structure

โครงสร้างนี้แยกหน้าที่ตาม production convention:

```text
ai-service/
├── app/
│   ├── main.py                  # FastAPI bootstrap only
│   ├── core/
│   │   └── settings.py          # env loading + settings model
│   ├── api/
│   │   └── routes/
│   │       ├── health.py        # GET /health
│   │       ├── chat.py          # POST /chat
│   │       └── embed.py         # POST /embed
│   ├── schemas/
│   │   ├── chat.py              # request/response models
│   │   └── embed.py
│   └── services/
│       ├── llm_service.py       # provider abstraction + chat/embed logic
│       └── memory_service.py    # Redis-backed chat memory (fallback in-memory)
├── requirements.txt
├── README.md
└── structure.md
```

## Design Rules
- `main.py` ห้ามมี business logic
- `routes/*` ทำแค่ validation + orchestration
- logic เรียก model/embedding ให้รวมที่ `services/`
- config ต้องมาจาก env ผ่าน `core/settings.py` เท่านั้น
- schema ต้องแยกไฟล์ตาม domain (`chat`, `embed`)

## Why this is production-friendly
- test ง่ายขึ้น (mock service ได้)
- maintainable เมื่อ endpoint เพิ่ม
- ลด coupling ระหว่าง HTTP layer กับ LLM layer
- รองรับ scale out / refactor เป็น microservice ง่าย
