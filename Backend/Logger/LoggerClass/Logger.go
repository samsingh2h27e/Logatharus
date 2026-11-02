package Logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	terminal *slog.Logger
	serviceName string
	isValid     bool
	conn *websocket.Conn
	token      string
}

type KeyValues struct {
    Key   string `json:"key"`
    Value any    `json:"value"`
}

type Event struct {
    Timestamp int64      `json:"timestamp"`
    TraceID   string     `json:"traceId"`
    Level     slog.Level `json:"level"`
    Message   string     `json:"message"`
    Map       []KeyValues `json:"map"`
}

type validationRequest struct {
	token string `json:"token" binding:"required"`
}

func NewLogger(serviceName string,token string) (*Logger , error) {
	if err := validateToken(token); err != nil {
		return nil,err
	}
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer "+ token)
	conn, resp, err := websocket.DefaultDialer.Dial(
        "ws://localhost:8080/Logger/Logs",
        headers,
    )
    if err != nil {
        fmt.Printf("WebSocket connection failed: %v\n", err)
        if resp != nil {
            fmt.Printf("HTTP status: %s\n", resp.Status)
        }
        return nil,err
    }

    fmt.Println("WebSocket connected successfully!")
	logger := &Logger{
		terminal : slog.New(slog.NewJSONHandler(os.Stdout,nil)),
		serviceName: serviceName,
		isValid:     false,
		conn: conn,
		token: 	token,
	}
	return logger,nil
}

func validateToken(token string) error {
	fmt.Println("Validating token:", token)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	reqBody := validationRequest{token: token}
	jsonData,err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}
	res,err := client.Post("http://localhost:8080/Logger/ValidateToken", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		return nil
	}else {
		return fmt.Errorf("token validation failed with status code: %d", res.StatusCode)
	}
}

func (L *Logger) Log(level slog.Level,ctx context.Context, message string, args ...any) {
	result := validateLogs(args...)
	if !result {
		log.Fatal("Invalid key-value pairs provided for logging. Keys must be strings.")
	}
	event := Event{
		Timestamp: time.Now().Unix(),
		TraceID:  GetTraceID(ctx),
		Level:     level,
		Message:   message,
	}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if !ok {
				continue
			}
			event.Map = append(event.Map, KeyValues{Key: key, Value: args[i+1]})
		}
	}
	fmt.Println("Logging event:", event)
	L.run(event)
	L.terminal.Log(ctx,level,message,args...)

}

func (L *Logger) Debug(ctx context.Context, message string, args ...any) {
	if !L.isValid {
		panic("Logger is not valid. Please check the token and connection.")
	}
	L.Log(slog.LevelDebug, ctx, message, args...)
}

func (L *Logger) Info(ctx context.Context, message string, args ...any) {
	if !L.isValid {
		panic("Logger is not valid. Please check the token and connection.")
	}
	L.Log(slog.LevelInfo, ctx, message, args...)
}

func (L *Logger) Warn(ctx context.Context, message string, args ...any) {
	if !L.isValid {
		panic("Logger is not valid. Please check the token and connection.")
	}
	L.Log(slog.LevelWarn, ctx, message, args...)
}

func (L *Logger) Error(ctx context.Context, message string, args ...any) {
	if !L.isValid {
		panic("Logger is not valid. Please check the token and connection.")
	}
	L.Log(slog.LevelError, ctx, message, args...)
}

func (L *Logger) run(event Event) {
    if L.conn == nil {
        return
    }
    if err := L.conn.WriteJSON(event); err != nil {
        slog.Error("Failed to send log event", "error", err)
    }
}

func (L *Logger) SendBatch(batch []Event, serviceName string) {
	// Implement the logic to send the batch of events to the desired destination (e.g., a logging service or database).
}

func validateLogs(args ...any) bool {
	for i := 0; i < len(args); i += 2 {
		_, ok := args[i].(string)
		if !ok {
			return false
		}
	}	
	return true
}

func GetTraceID(ctx context.Context) string {
    spanContext := trace.SpanContextFromContext(ctx)

    if !spanContext.IsValid() {
        return ""
    }

    return spanContext.TraceID().String()
}