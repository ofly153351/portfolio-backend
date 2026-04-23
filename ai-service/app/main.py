from fastapi import FastAPI

from app.api.routes.admin_memory import router as admin_memory_router
from app.api.routes.chat import router as chat_router
from app.api.routes.embed import router as embed_router
from app.api.routes.health import router as health_router


def create_app() -> FastAPI:
    app = FastAPI(title="portfolio-ai-service", version="1.0.0")
    app.include_router(health_router)
    app.include_router(chat_router)
    app.include_router(embed_router)
    app.include_router(admin_memory_router)
    return app


app = create_app()
