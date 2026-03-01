# AI Voice Agent - Day 1

## What this is
A Day-1 skeleton for an AI assistant that:
- streams responses using Server-Sent Events (SSE)
- calls a tool endpoint (`/tools/get_calendar_slots`) and responds using tool data

## Python server
### Run
cd ai-voice-agent/python
source .venv/bin/activate
uvicorn app.main:app --reload --port 8000

### Test
curl http://127.0.0.1:8000/health

curl -N -X POST http://127.0.0.1:8000/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"message":"hi"}'

curl -N -X POST http://127.0.0.1:8000/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"message":"schedule an appointment"}'
