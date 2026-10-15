package handler

import (
	"github.com/gin-gonic/gin"
)

type LogRequest struct {
	ServiceName string     `json:"serviceName" binding:"required"`
	Batch       []LogBatch `json:"batch" binding:"required"`
}

type LogBatch struct {
	Timestamp int64      `json:"timestamp"`
	TraceID   string     `json:"traceId"`
	Level     string     `json:"level"`
	Message   string     `json:"message"`
	Map       []KeyValue `json:"map"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}


func Log(ctx *gin.Context) {
	// Implement your logging logic here
	var requestBody LogRequest
	err := ctx.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
}