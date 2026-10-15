package Logger

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	terminal *slog.Logger
	ServiceName string
	channel     chan Event
	isValid     bool
}

type KeyValues struct {
	key  string
	value any
}

type Event struct {
	Timestamp   int64
	traceID	    string
	Level       slog.Level
	Message     string
	Map 	    []KeyValues
}

func NewLogger(serviceName string) *Logger {
	logger := &Logger{
		terminal : slog.New(slog.NewJSONHandler(os.Stdout,nil)),
		ServiceName: serviceName,
		channel : make(chan Event, 1000),
		isValid:     false,
	}
	go logger.run()
	return logger
}

func (L *Logger) ValidateToken(token string) error {
	return nil
}

func (L *Logger) Log(level slog.Level,ctx context.Context, message string, args ...any) {
	result := validateLogs(args...)
	if !result {
		log.Fatal("Invalid key-value pairs provided for logging. Keys must be strings.")
	}
	event := Event{
		Timestamp: time.Now().Unix(),
		traceID:  GetTraceID(ctx),
		Level:     level,
		Message:   message,
	}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if !ok {
				continue
			}
			event.Map = append(event.Map, KeyValues{key: key, value: args[i+1]})
		}
	}
	L.channel <- event
	L.terminal.Log(ctx,level,message,args...)

}

func (L *Logger) Debug(ctx context.Context, message string, args ...any) {
	L.Log(slog.LevelDebug, ctx, message, args...)
}

func (L *Logger) Info(ctx context.Context, message string, args ...any) {
	L.Log(slog.LevelInfo, ctx, message, args...)
}

func (L *Logger) Warn(ctx context.Context, message string, args ...any) {
	L.Log(slog.LevelWarn, ctx, message, args...)
}

func (L *Logger) Error(ctx context.Context, message string, args ...any) {
	L.Log(slog.LevelError, ctx, message, args...)
}

func (L *Logger) run() {
	batch := make([]Event, 0, 100)
	for {
		event := <-L.channel
		batch = append(batch, event)
		if len(batch) >= 100 {
			L.SendBatch(batch,L.ServiceName)
			batch = make([]Event, 0,100)	
		}
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