from fastapi import APIRouter, HTTPException

from app.core.settings import load_settings
from app.schemas.chat import ChatRequest, ChatResponse
from app.services.llm_service import LLMService

router = APIRouter(tags=["chat"])


def _service() -> LLMService:
    return LLMService(load_settings())


@router.post("/chat", response_model=ChatResponse)
def chat(req: ChatRequest) -> ChatResponse:
    if not req.message.strip():
        raise HTTPException(status_code=400, detail="message is required")

    result = _service().ask(req.message, req.session_id)
    return ChatResponse(
        answer=result.answer,
        sources=[],
        session_id=req.session_id,
        provider=result.provider,
        usage={
            "prompt_tokens": result.usage.prompt_tokens,
            "completion_tokens": result.usage.completion_tokens,
            "total_tokens": result.usage.total_tokens,
        },
    )
