import json
from dataclasses import dataclass
from threading import RLock
from typing import Optional, Protocol

import redis

from app.core.settings import Settings


@dataclass(frozen=True)
class ChatTurn:
    role: str
    content: str


class MemoryStore(Protocol):
    def get_history(self, session_id: str) -> list[ChatTurn]:
        ...

    def append_turn(self, session_id: str, role: str, content: str, max_turns: int) -> None:
        ...

    def clear(self, session_id: str) -> None:
        ...


class InMemoryStore:
    def __init__(self) -> None:
        self._lock = RLock()
        self._sessions: dict[str, list[ChatTurn]] = {}

    def get_history(self, session_id: str) -> list[ChatTurn]:
        if not session_id:
            return []
        with self._lock:
            return list(self._sessions.get(session_id, []))

    def append_turn(self, session_id: str, role: str, content: str, max_turns: int) -> None:
        if not session_id or not content:
            return
        safe_max_turns = max(max_turns, 1)
        with self._lock:
            turns = self._sessions.setdefault(session_id, [])
            turns.append(ChatTurn(role=role, content=content))
            if len(turns) > safe_max_turns * 2:
                del turns[: len(turns) - safe_max_turns*2]

    def clear(self, session_id: str) -> None:
        if not session_id:
            return
        with self._lock:
            self._sessions.pop(session_id, None)


class RedisMemoryStore:
    def __init__(self, settings: Settings) -> None:
        self.settings = settings
        self.client = redis.Redis.from_url(
            settings.redis_url,
            decode_responses=True,
            socket_timeout=settings.redis_timeout_seconds,
            socket_connect_timeout=settings.redis_timeout_seconds,
        )

    def _key(self, session_id: str) -> str:
        return f"{self.settings.redis_key_prefix}{session_id}"

    def get_history(self, session_id: str) -> list[ChatTurn]:
        if not session_id:
            return []
        rows = self.client.lrange(self._key(session_id), 0, -1)
        turns: list[ChatTurn] = []
        for row in rows:
            try:
                obj = json.loads(row)
            except json.JSONDecodeError:
                continue
            role = str(obj.get("role", "")).strip()
            content = str(obj.get("content", "")).strip()
            if role in {"user", "assistant"} and content:
                turns.append(ChatTurn(role=role, content=content))
        return turns

    def append_turn(self, session_id: str, role: str, content: str, max_turns: int) -> None:
        if not session_id or not content:
            return
        safe_max_turns = max(max_turns, 1)
        key = self._key(session_id)
        payload = json.dumps({"role": role, "content": content}, ensure_ascii=False)
        pipe = self.client.pipeline(transaction=False)
        pipe.rpush(key, payload)
        pipe.ltrim(key, -safe_max_turns * 2, -1)
        pipe.execute()

    def clear(self, session_id: str) -> None:
        if not session_id:
            return
        self.client.delete(self._key(session_id))


class SafeMemoryStore:
    def __init__(self, primary: MemoryStore, fallback: MemoryStore) -> None:
        self.primary = primary
        self.fallback = fallback

    def get_history(self, session_id: str) -> list[ChatTurn]:
        try:
            return self.primary.get_history(session_id)
        except Exception:
            return self.fallback.get_history(session_id)

    def append_turn(self, session_id: str, role: str, content: str, max_turns: int) -> None:
        try:
            self.primary.append_turn(session_id, role, content, max_turns)
            return
        except Exception:
            pass
        self.fallback.append_turn(session_id, role, content, max_turns)

    def clear(self, session_id: str) -> None:
        try:
            self.primary.clear(session_id)
            return
        except Exception:
            pass
        self.fallback.clear(session_id)


_store: Optional[MemoryStore] = None
_store_lock = RLock()


def get_memory_store(settings: Settings) -> MemoryStore:
    global _store
    with _store_lock:
        if _store is not None:
            return _store

        fallback = InMemoryStore()
        if settings.memory_store == "memory":
            _store = fallback
            return _store

        try:
            primary = RedisMemoryStore(settings)
            primary.client.ping()
            _store = SafeMemoryStore(primary=primary, fallback=fallback)
            return _store
        except Exception:
            _store = fallback
            return _store
