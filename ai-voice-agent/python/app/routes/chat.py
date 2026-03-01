import asyncio
import json
from datetime import datetime, timedelta

import httpx
from fastapi import APIRouter
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

router = APIRouter(prefix="/chat")

class ChatRequest(BaseModel):
    message: str

def _sse(data: dict) -> str:
    return f"data: {json.dumps(data)}\n\n"

async def _stream_words(text: str):
    for word in text.split(" "):
        yield _sse({"token": word + " "})
        await asyncio.sleep(0.05)
    yield _sse({"done": True})

@router.post("/stream")
async def chat_stream(req: ChatRequest):
    msg = req.message.lower()

    async def event_gen():
        # Normal path
        if "appointment" not in msg and "schedule" not in msg:
            async for chunk in _stream_words(
                "Hello! I am your AI assistant. Ask me to schedule an appointment."
            ):
                yield chunk
            return

        # Tool calling path (simulated)
        yield _sse({"event": "tool_call", "tool": "get_calendar_slots"})

        tomorrow = (datetime.utcnow() + timedelta(days=1)).date().isoformat()

        async with httpx.AsyncClient(timeout=10.0) as client:
            resp = await client.post(
                "http://127.0.0.1:8000/tools/get_calendar_slots",
                json={"date": tomorrow},
            )
            resp.raise_for_status()
            tool_data = resp.json()

        yield _sse({"event": "tool_result", "tool": "get_calendar_slots", "result": tool_data})

        slots = tool_data.get("slots", [])
        slot_text = ", ".join(slots) if slots else "No slots available"

        async for chunk in _stream_words(
            f"Sure. Available slots for {tool_data['date']} are: {slot_text}. Which one do you prefer?"
        ):
            yield chunk

    return StreamingResponse(event_gen(), media_type="text/event-stream")
