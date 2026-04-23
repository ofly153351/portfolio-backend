# Local AI Setup (Ollama)

## Recommended small models
- Chat model (small): `llama3.2:1b`
- Embedding model (small): `nomic-embed-text`

> ถ้าต้องการไทยดีขึ้นเล็กน้อย (แต่หนักขึ้น): ใช้ `qwen2.5:1.5b`

## 1) Pull models
```bash
ollama pull llama3.2:1b
ollama pull nomic-embed-text
```

## 2) Start Ollama
```bash
ollama serve
```

## 3) Configure backend env (`.env`)
```env
# AI service -> Ollama (OpenAI-compatible endpoint)
OPENAI_API_KEY=ollama
OPENAI_BASE_URL=http://localhost:11434/v1
OPENAI_CHAT_MODEL=llama3.2:1b
OPENAI_EMBED_MODEL=nomic-embed-text

# optional
OPENROUTER_HTTP_REFERER=
OPENROUTER_X_TITLE=portfolio-backend
```

## 4) Run AI service
```bash
cd ai-service
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

## 5) Test endpoints

### Chat (through Go API)
```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"สวัสดี ช่วยสรุป portfolio ให้หน่อย","session_id":"local-1"}'
```

### Embed (direct to AI service)
```bash
curl -X POST http://localhost:8000/embed \
  -H "Content-Type: application/json" \
  -d '{"id":"cnt_001","title":"POS","content":"ระบบขายหน้าร้าน","type":"project"}'
```

## Notes
- ถ้า RAM น้อยมาก ให้เริ่มจาก `llama3.2:1b` ก่อน
- ถ้าตอบช้า ให้ปิด model อื่นในเครื่อง
