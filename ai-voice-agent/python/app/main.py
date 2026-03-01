from fastapi import FastAPI

from app.routes.health import router as health_router
from app.routes.tools import router as tools_router
from app.routes.chat import router as chat_router

app = FastAPI(title="ai-voice-agent-python")

app.include_router(health_router)
app.include_router(tools_router)
app.include_router(chat_router)
