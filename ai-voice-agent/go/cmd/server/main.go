package main

import (
	"ai-voice-agent-go/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", handlers.Health)

	tools := r.Group("/tools")
	tools.POST("/get_calendar_slots", handlers.GetCalendarSlots)

	chat := r.Group("/chat")
	chat.POST("/stream", handlers.ChatStream)

	// Go server runs on 9000 so it doesn't conflict with Python 8000
	r.Run(":9000")
}
