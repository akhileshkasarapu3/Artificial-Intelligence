package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func writeSSE(w *bufio.Writer, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := w.WriteString("data: " + string(b) + "\n\n"); err != nil {
		return err
	}
	return w.Flush()
}

func streamWords(w *bufio.Writer, text string) error {
	words := strings.Split(text, " ")
	for _, word := range words {
		if err := writeSSE(w, gin.H{"token": word + " "}); err != nil {
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
	return writeSSE(w, gin.H{"done": true})
}

func ChatStream(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	msg := strings.ToLower(req.Message)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	w := bufio.NewWriter(c.Writer)

	// Normal path
	if !strings.Contains(msg, "appointment") && !strings.Contains(msg, "schedule") {
		_ = streamWords(w, "Hello! I am your AI assistant. Ask me to schedule an appointment.")
		return
	}

	// Tool-calling path (simulated)
	_ = writeSSE(w, gin.H{"event": "tool_call", "tool": "get_calendar_slots"})

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	body, _ := json.Marshal(gin.H{"date": tomorrow})

	httpResp, err := http.Post(
		"http://127.0.0.1:9000/tools/get_calendar_slots",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		_ = streamWords(w, "Sorry, I could not fetch calendar slots right now.")
		return
	}
	defer httpResp.Body.Close()

	var toolResult map[string]any
	if err := json.NewDecoder(httpResp.Body).Decode(&toolResult); err != nil {
		_ = streamWords(w, "Sorry, I received an invalid tool response.")
		return
	}

	_ = writeSSE(w, gin.H{"event": "tool_result", "tool": "get_calendar_slots", "result": toolResult})

	slotsAny, _ := toolResult["slots"].([]any)
	slots := make([]string, 0, len(slotsAny))
	for _, s := range slotsAny {
		if v, ok := s.(string); ok {
			slots = append(slots, v)
		}
	}

	slotText := "No slots available"
	if len(slots) > 0 {
		slotText = strings.Join(slots, ", ")
	}

	dateStr, _ := toolResult["date"].(string)
	_ = streamWords(w, "Sure. Available slots for "+dateStr+" are: "+slotText+". Which one do you prefer?")
}
