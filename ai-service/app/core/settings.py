import os
from pathlib import Path
from dataclasses import dataclass

from dotenv import load_dotenv


def _load_env() -> None:
    current_file = Path(__file__).resolve()
    project_root = current_file.parents[3]
    env_path = project_root / ".env"
    if env_path.exists():
        load_dotenv(env_path, override=False)


@dataclass(frozen=True)
class Settings:
    api_service: str
    openai_api_key: str
    openai_base_url: str
    openai_chat_model: str
    openai_embed_model: str
    openai_timeout_seconds: int
    ollama_base_url: str
    ollama_model: str
    ollama_embed_model: str
    ollama_timeout_seconds: int
    openrouter_http_referer: str
    openrouter_x_title: str
    system_prompt: str
    memory_max_turns: int
    memory_store: str
    redis_url: str
    redis_key_prefix: str
    redis_timeout_seconds: int


def load_settings() -> Settings:
    _load_env()
    api_service = os.getenv("API_SERVICE", "openai").strip().lower()
    if api_service not in {"openai", "ollama"}:
        api_service = "openai"

    openai_timeout = 60
    try:
        openai_timeout = int(os.getenv("OPENAI_TIMEOUT_SECONDS", "60").strip())
    except ValueError:
        openai_timeout = 60
    if openai_timeout <= 0:
        openai_timeout = 60

    ollama_timeout = 300
    try:
        ollama_timeout = int(os.getenv("OLLAMA_TIMEOUT_SECONDS", "300").strip())
    except ValueError:
        ollama_timeout = 300
    if ollama_timeout <= 0:
        ollama_timeout = 300

    memory_max_turns = 12
    try:
        memory_max_turns = int(os.getenv("MEMORY_MAX_TURNS", "12").strip())
    except ValueError:
        memory_max_turns = 12
    if memory_max_turns <= 0:
        memory_max_turns = 12

    memory_store = os.getenv("MEMORY_STORE", "redis").strip().lower()
    if memory_store not in {"redis", "memory"}:
        memory_store = "redis"

    redis_timeout = 5
    try:
        redis_timeout = int(os.getenv("REDIS_TIMEOUT_SECONDS", "5").strip())
    except ValueError:
        redis_timeout = 5
    if redis_timeout <= 0:
        redis_timeout = 5

    return Settings(
        api_service=api_service,
        openai_api_key=os.getenv("OPENAI_API_KEY", "").strip(),
        openai_base_url=os.getenv("OPENAI_BASE_URL", "https://api.openai.com/v1").strip(),
        openai_chat_model=os.getenv("OPENAI_CHAT_MODEL", "gpt-4o-mini").strip(),
        openai_embed_model=os.getenv("OPENAI_EMBED_MODEL", "text-embedding-3-small").strip(),
        openai_timeout_seconds=openai_timeout,
        ollama_base_url=os.getenv("OLLAMA_BASE_URL", "http://localhost:11434").strip(),
        ollama_model=os.getenv("OLLAMA_MODEL", "gemma3:1b").strip(),
        ollama_embed_model=os.getenv("OLLAMA_EMBED_MODEL", "nomic-embed-text").strip(),
        ollama_timeout_seconds=ollama_timeout,
        openrouter_http_referer=os.getenv("OPENROUTER_HTTP_REFERER", "").strip(),
        openrouter_x_title=os.getenv("OPENROUTER_X_TITLE", "portfolio-ai-service").strip(),
        system_prompt=os.getenv(
            "AI_SYSTEM_PROMPT",
            "You are Peerapat. Style: practical, opinionated, dev mindset.",
        ).strip(),
        memory_max_turns=memory_max_turns,
        memory_store=memory_store,
        redis_url=os.getenv("REDIS_URL", "redis://:password@localhost:6379/0").strip(),
        redis_key_prefix=os.getenv("REDIS_KEY_PREFIX", "portfolio:chat:").strip() or "portfolio:chat:",
        redis_timeout_seconds=redis_timeout,
    )
