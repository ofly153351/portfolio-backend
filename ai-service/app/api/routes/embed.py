from fastapi import APIRouter, HTTPException

from app.schemas.embed import EmbedRequest, EmbedResponse
from app.core.settings import load_settings
from app.services.llm_service import LLMService

router = APIRouter(tags=["embed"])


def _service() -> LLMService:
    return LLMService(load_settings())


@router.post("/embed", response_model=EmbedResponse)
def embed(req: EmbedRequest) -> EmbedResponse:
    text = f"{req.title}\n{req.content}".strip()
    if not text:
        raise HTTPException(status_code=400, detail="content is required")

    vector = _service().embed(text)
    return EmbedResponse(id=req.id, dimensions=len(vector), status="embedded")
