from fastapi import APIRouter, HTTPException

from app.core.settings import load_settings
from app.services.memory_service import get_memory_store

router = APIRouter(tags=["admin-memory"])


@router.get("/admin/chat-memory/{session_id}")
def get_chat_memory(session_id: str) -> dict:
    sid = session_id.strip()
    if not sid:
        raise HTTPException(status_code=400, detail="session_id is required")

    settings = load_settings()
    store = get_memory_store(settings)
    history = store.get_history(sid)
    turns = [{"role": turn.role, "content": turn.content} for turn in history]
    return {
        "session_id": sid,
        "memory_store": settings.memory_store,
        "count": len(turns),
        "turns": turns,
    }


@router.delete("/admin/chat-memory/{session_id}")
def delete_chat_memory(session_id: str) -> dict:
    sid = session_id.strip()
    if not sid:
        raise HTTPException(status_code=400, detail="session_id is required")

    settings = load_settings()
    store = get_memory_store(settings)
    before = len(store.get_history(sid))
    store.clear(sid)
    return {
        "session_id": sid,
        "memory_store": settings.memory_store,
        "deleted": before > 0,
        "cleared_count": before,
    }
