package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CalendarSlotsRequest struct {
	Date string `json:"date"`
}

func GetCalendarSlots(c *gin.Context) {
	var req CalendarSlotsRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":  req.Date,
		"slots": []string{"10:00", "11:30", "15:00"},
	})
}
