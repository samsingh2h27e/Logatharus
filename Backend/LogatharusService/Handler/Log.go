package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Development only; restrict origins in production.
	},
}

func Log(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade to WebSocket", "error", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			slog.Info("WebSocket connection closed", "error", err)
			return
		}

		slog.Info(
			"Received message",
			"messageType", messageType,
			"message", string(message),
		)
	}
}