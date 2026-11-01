package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Event struct {
    Timestamp int64      `json:"timestamp"`
    TraceID   string     `json:"traceId"`
    Level     slog.Level `json:"level"`
    Message   string     `json:"message"`
    Map       []KeyValues `json:"map"`
}

type KeyValues struct {
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
	auth := ctx.GetHeader("Authorization")
	user := getUserFromToken(auth)
	storage := make(map[string][]Event)
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
		if _,ok := storage[user]; !ok {
			storage[user] = []Event{}
		}
		storage[user] = append(storage[user], convertToStruct(message))
	}
}

func getUserFromToken(token string) string {
	// Placeholder for token validation logic.
	// In a real implementation, validate the token and extract user information.
	return "user123"
}

func convertToStruct(message []byte) Event {
	var event Event
	err := json.Unmarshal(message, &event)
	if err != nil {
		slog.Error("Failed to unmarshal message", "error", err)
	}
	return event
}