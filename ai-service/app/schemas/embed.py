from pydantic import BaseModel


class EmbedRequest(BaseModel):
    id: str
    title: str = ""
    content: str
    type: str = "content"


class EmbedResponse(BaseModel):
    id: str
    dimensions: int
    status: str
