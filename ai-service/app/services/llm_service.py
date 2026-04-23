import json
from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Dict, Iterator
from urllib import error as urlerror
from urllib import request as urlrequest

from fastapi import HTTPException
from langchain_openai import OpenAIEmbeddings

from app.core.settings import Settings
from app.services.memory_service import ChatTurn, MemoryStore, get_memory_store


@dataclass(frozen=True)
class ChatUsage:
    prompt_tokens: int = 0
    completion_tokens: int = 0
    total_tokens: int = 0


@dataclass(frozen=True)
class ChatResult:
    answer: str
    provider: str
    usage: ChatUsage


class ProviderService(ABC):
    @abstractmethod
    def chat(self, prompt: str, history: list[ChatTurn]) -> ChatResult:
        raise NotImplementedError


class OpenAIService(ProviderService):
    def __init__(self, settings: Settings) -> None:
        self.settings = settings

    def _headers(self) -> Dict[str, str]:
        headers: Dict[str, str] = {
            "Authorization": f"Bearer {self.settings.openai_api_key}",
            "Content-Type": "application/json",
        }
        if self.settings.openrouter_http_referer:
            headers["HTTP-Referer"] = self.settings.openrouter_http_referer
        if self.settings.openrouter_x_title:
            headers["X-Title"] = self.settings.openrouter_x_title
        return headers

    def chat(self, prompt: str, history: list[ChatTurn]) -> ChatResult:
        if not self.settings.openai_api_key:
            raise HTTPException(status_code=500, detail="OPENAI_API_KEY is missing")

        url = self.settings.openai_base_url.rstrip("/") + "/chat/completions"
        payload = {
            "model": self.settings.openai_chat_model,
            "messages": [
                {"role": "system", "content": self.settings.system_prompt},
            ],
        }
        for turn in history:
            if turn.role in {"user", "assistant"} and turn.content:
                payload["messages"].append({"role": turn.role, "content": turn.content})
        payload["messages"].append({"role": "user", "content": prompt})
        body = json.dumps(payload).encode("utf-8")
        req = urlrequest.Request(url, data=body, headers=self._headers(), method="POST")

        try:
            with urlrequest.urlopen(req, timeout=self.settings.openai_timeout_seconds) as resp:
                raw = resp.read().decode("utf-8")
        except urlerror.HTTPError as err:
            detail = err.read().decode("utf-8", errors="replace")
            raise HTTPException(status_code=err.code, detail=detail) from err
        except Exception as err:  # noqa: BLE001
            raise HTTPException(status_code=500, detail=f"Connection error: {err}") from err

        try:
            data = json.loads(raw)
            usage = data.get("usage", {})
            prompt_tokens = int(usage.get("prompt_tokens", 0) or 0)
            completion_tokens = int(usage.get("completion_tokens", 0) or 0)
            total_tokens = int(
                usage.get("total_tokens", prompt_tokens + completion_tokens) or 0
            )
            return ChatResult(
                answer=data["choices"][0]["message"]["content"],
                provider="openai",
                usage=ChatUsage(
                    prompt_tokens=prompt_tokens,
                    completion_tokens=completion_tokens,
                    total_tokens=total_tokens,
                ),
            )
        except Exception as err:  # noqa: BLE001
            raise HTTPException(status_code=500, detail=f"Invalid OpenAI response: {err}") from err


class OllamaService(ProviderService):
    def __init__(self, settings: Settings) -> None:
        self.settings = settings

    def chat(self, prompt: str, history: list[ChatTurn]) -> ChatResult:
        url = self.settings.ollama_base_url.rstrip("/") + "/api/generate"
        history_text = "\n".join(
            f"{'User' if turn.role == 'user' else 'Assistant'}: {turn.content}"
            for turn in history
            if turn.content
        )
        history_block = f"{history_text}\n" if history_text else ""
        full_prompt = (
            f"{self.settings.system_prompt}\n\n"
            f"{history_block}"
            f"User: {prompt}\n"
            f"Assistant:"
        )
        payload = {
            "model": self.settings.ollama_model,
            "prompt": full_prompt,
            "stream": False,
        }
        body = json.dumps(payload).encode("utf-8")
        req = urlrequest.Request(
            url,
            data=body,
            headers={"Content-Type": "application/json"},
            method="POST",
        )

        try:
            with urlrequest.urlopen(req, timeout=self.settings.ollama_timeout_seconds) as resp:
                raw = resp.read().decode("utf-8")
        except urlerror.HTTPError as err:
            detail = err.read().decode("utf-8", errors="replace")
            raise HTTPException(status_code=err.code, detail=detail) from err
        except Exception as err:  # noqa: BLE001
            raise HTTPException(status_code=500, detail=f"Connection error: {err}") from err

        try:
            data = json.loads(raw)
            prompt_tokens = int(data.get("prompt_eval_count", 0) or 0)
            completion_tokens = int(data.get("eval_count", 0) or 0)
            return ChatResult(
                answer=data.get("response", ""),
                provider="ollama",
                usage=ChatUsage(
                    prompt_tokens=prompt_tokens,
                    completion_tokens=completion_tokens,
                    total_tokens=prompt_tokens + completion_tokens,
                ),
            )
        except Exception as err:  # noqa: BLE001
            raise HTTPException(status_code=500, detail=f"Invalid Ollama response: {err}") from err


def get_llm_service(settings: Settings) -> ProviderService:
    if settings.api_service == "ollama":
        return OllamaService(settings)
    return OpenAIService(settings)


class LLMService:
    def __init__(self, settings: Settings) -> None:
        self.settings = settings
        self.provider = get_llm_service(settings)
        self.memory_store: MemoryStore = get_memory_store(settings)

    def ask(self, message: str, session_id: str) -> ChatResult:
        sid = (session_id or "session-default").strip() or "session-default"
        history = self.memory_store.get_history(sid)
        result = self.provider.chat(message, history)
        self.memory_store.append_turn(sid, "user", message, self.settings.memory_max_turns)
        self.memory_store.append_turn(sid, "assistant", result.answer, self.settings.memory_max_turns)
        return result

    def stream(self, message: str) -> Iterator[str]:
        result = self.provider.chat(message, [])
        for token in result.answer.split():
            yield token + " "

    def embed(self, text: str) -> list[float]:
        if self.settings.api_service == "ollama":
            return self._embed_with_ollama(text)
        return self._embed_with_openai(text)

    def _embed_with_openai(self, text: str) -> list[float]:
        if not self.settings.openai_api_key:
            raise HTTPException(status_code=500, detail="OPENAI_API_KEY is missing")
        embeddings = OpenAIEmbeddings(
            api_key=self.settings.openai_api_key,
            base_url=self.settings.openai_base_url,
            model=self.settings.openai_embed_model,
        )
        try:
            return embeddings.embed_query(text)
        except Exception as err:  # noqa: BLE001
            raise HTTPException(status_code=500, detail=str(err)) from err

    def _embed_with_ollama(self, text: str) -> list[float]:
        url = self.settings.ollama_base_url.rstrip("/") + "/api/embeddings"
        payload = {
            "model": self.settings.ollama_embed_model or self.settings.ollama_model,
            "prompt": text,
        }
        body = json.dumps(payload).encode("utf-8")
        req = urlrequest.Request(
            url,
            data=body,
            headers={"Content-Type": "application/json"},
            method="POST",
        )

        try:
            with urlrequest.urlopen(req, timeout=self.settings.ollama_timeout_seconds) as resp:
                raw = resp.read().decode("utf-8")
        except urlerror.HTTPError as err:
            detail = err.read().decode("utf-8", errors="replace")
            raise HTTPException(status_code=err.code, detail=detail) from err
        except Exception as err:  # noqa: BLE001
            raise HTTPException(status_code=500, detail=f"Connection error: {err}") from err

        try:
            data = json.loads(raw)
            embedding = data.get("embedding")
            if not isinstance(embedding, list):
                raise ValueError("missing embedding field")
            return [float(x) for x in embedding]
        except Exception as err:  # noqa: BLE001
            raise HTTPException(status_code=500, detail=f"Invalid Ollama embedding response: {err}") from err
