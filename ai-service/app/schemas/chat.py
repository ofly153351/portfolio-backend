from typing import List

from pydantic import BaseModel


class ChatRequest(BaseModel):
    message: str
    session_id: str = "session-default"
    top_k: int = 5
    lang: str = "th"


class ChatSource(BaseModel):
    id: str
    title: str
    score: float


class ChatUsage(BaseModel):
    prompt_tokens: int = 0
    completion_tokens: int = 0
    total_tokens: int = 0


class ChatResponse(BaseModel):
    answer: str
    sources: List[ChatSource]
    session_id: str
    provider: str = ""
    usage: ChatUsage = ChatUsage()


class StreamEvent(BaseModel):
    type: str
    token: str = ""
    message: str = ""
    error: str = ""
    session_id: str = ""
