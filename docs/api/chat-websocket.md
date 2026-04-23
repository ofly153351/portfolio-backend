# WebSocket Chat Streaming (Only Streaming Channel)

## Endpoint
- URL: `ws://localhost:8080/api/chat/ws`
- Protocol: WebSocket

## Client -> Server message
```json
{
  "message": "ช่วยสรุป portfolio ให้หน่อย",
  "session_id": "sess_ws_001",
  "top_k": 5,
  "lang": "th"
}
```

## Server -> Client events

### status
```json
{ "type": "status", "message": "streaming" }
```

### token
```json
{ "type": "token", "token": "..." }
```

### done
```json
{
  "type": "done",
  "session_id": "sess_ws_001",
  "provider": "ollama",
  "usage": {
    "prompt_tokens": 142,
    "completion_tokens": 96,
    "total_tokens": 238
  }
}
```

### error
```json
{ "type": "error", "error": "..." }
```

## Notes
- Streaming is only on WebSocket.
- Go API uses AI service `POST /chat` behind the scenes, then emits token events to client.
- Frontend should concat `token` events to render realtime response.
- Token usage จริงของ provider ต่อรอบ อ่านจาก event `type=done` -> `usage`.
- บุคลิก/กฎการตอบของโมเดลอ่านจาก env `AI_SYSTEM_PROMPT` (root `.env`).
- Memory ทำงานตาม `session_id`:
  - ส่ง `session_id` เดิม -> ต่อบทสนทนาเดิม
  - ส่ง `session_id` ใหม่ -> เริ่ม memory ใหม่
- memory backend ปัจจุบันเก็บใน Redis (persist ข้าม restart ของ ai-service ได้)

## Next.js Example (Client Component)
```tsx
"use client";

import { useRef, useState } from "react";

type WsEvent = {
  type: "status" | "token" | "done" | "error";
  token?: string;
  message?: string;
  error?: string;
  session_id?: string;
  provider?: string;
  usage?: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
};

export default function ChatWS() {
  const wsRef = useRef<WebSocket | null>(null);
  const [text, setText] = useState("");
  const [reply, setReply] = useState("");
  const [status, setStatus] = useState("idle");
  const [usage, setUsage] = useState<WsEvent["usage"]>();

  const connect = () => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) return;
    const ws = new WebSocket("ws://localhost:8080/api/chat/ws");

    ws.onopen = () => setStatus("connected");
    ws.onclose = () => setStatus("closed");
    ws.onerror = () => setStatus("error");

    ws.onmessage = (event) => {
      const data: WsEvent = JSON.parse(event.data);
      if (data.type === "status") setStatus(data.message ?? "streaming");
      if (data.type === "token") setReply((prev) => prev + (data.token ?? ""));
      if (data.type === "done") {
        setStatus("done");
        setUsage(data.usage);
      }
      if (data.type === "error") setStatus(data.error ?? "error");
    };

    wsRef.current = ws;
  };

  const send = () => {
    const ws = wsRef.current;
    if (!ws || ws.readyState !== WebSocket.OPEN) return;
    setReply("");
    setUsage(undefined);
    ws.send(
      JSON.stringify({
        message: text,
        session_id: "nextjs-ws-1",
        top_k: 5,
        lang: "th",
      })
    );
  };

  return (
    <div>
      <button onClick={connect}>Connect</button>
      <input value={text} onChange={(e) => setText(e.target.value)} />
      <button onClick={send}>Send</button>
      <p>Status: {status}</p>
      {usage && (
        <p>
          Usage: prompt={usage.prompt_tokens}, completion={usage.completion_tokens},
          total={usage.total_tokens}
        </p>
      )}
      <pre>{reply}</pre>
    </div>
  );
}
```
