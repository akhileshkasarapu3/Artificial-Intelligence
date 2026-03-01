# AI Voice Agent — Day 1

## What this is
A Day-1 skeleton for an AI assistant that:
- Streams responses using Server-Sent Events (SSE)
- Calls a tool endpoint (`/tools/get_calendar_slots`) and responds using tool data
- Works in both **Python (FastAPI)** and **Go (Gin)**

---

## Folder structure
```text
ai-voice-agent/
  python/
    app/
      main.py
      routes/
        health.py
        tools.py
        chat.py
    requirements.txt
    .venv/
  go/
    go.mod
    cmd/server/main.go
    internal/handlers/
      health.go
      tools.go
      chat.go
```

---

## Python server (FastAPI)

### Prerequisites
- Python **3.12+** recommended

### Setup (first time only)
```bash
cd ai-voice-agent/python
python3 -m venv .venv
source .venv/bin/activate
pip install --upgrade pip
pip install -r requirements.txt
```

### Run
```bash
cd ai-voice-agent/python
source .venv/bin/activate
uvicorn app.main:app --reload --port 8000
```

### Test
Health:
```bash
curl http://127.0.0.1:8000/health
```

Streaming:
```bash
curl -N -X POST http://127.0.0.1:8000/chat/stream   -H "Content-Type: application/json"   -d '{"message":"hi"}'
```

Tool calling + streaming:
```bash
curl -N -X POST http://127.0.0.1:8000/chat/stream   -H "Content-Type: application/json"   -d '{"message":"schedule an appointment"}'
```

### Notes
- SSE works because the response `Content-Type` is `text/event-stream`.
- Each event ends with `\n\n` (blank line) so clients can parse event boundaries.
- `/tools/*` endpoints are deterministic “functions” (mocked for Day 1).
- `/chat/*` orchestrates conversation + tool calling + streaming.

---

## Go server (Gin)

### Prerequisites
- Go **1.22+** recommended

### Setup (first time only)
```bash
cd ai-voice-agent/go
go get github.com/gin-gonic/gin@v1.10.0
go mod tidy
```

### Run
```bash
cd ai-voice-agent/go
go run ./cmd/server
```

> Go server runs on **port 9000** so it doesn’t conflict with Python (8000).

### Test
Health:
```bash
curl http://127.0.0.1:9000/health
```

Streaming:
```bash
curl -N -X POST http://127.0.0.1:9000/chat/stream   -H "Content-Type: application/json"   -d '{"message":"hi"}'
```

Tool calling + streaming:
```bash
curl -N -X POST http://127.0.0.1:9000/chat/stream   -H "Content-Type: application/json"   -d '{"message":"schedule an appointment"}'
```

### Notes
- Go streams SSE by writing `data: <json>\n\n` repeatedly and flushing.
- Tool calling is simulated by the chat handler calling:
  - `POST http://127.0.0.1:9000/tools/get_calendar_slots`

---

## Day 1 completion criteria
You’re done when:
- Python streams tokens and shows `tool_call` → `tool_result`
- Go streams tokens and shows `tool_call` → `tool_result`
- Both servers return `{"status":"ok"}` for `/health`
- You can run both services from scratch using this README
