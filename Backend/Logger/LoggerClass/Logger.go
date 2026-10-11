package Logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	terminal *slog.Logger
	ServiceName string
	channel     chan Event
	isValid     bool
}

type Event struct {

}

func NewLogger(serviceName string) *Logger {
	logger := &Logger{
		terminal : slog.New(slog.NewJSONHandler(os.Stdout,nil)),
		ServiceName: serviceName,
		channel : make(chan Event, 1000),
		isValid:     false,
	}
	go logger.Run()
	return logger
}

func (L *Logger) ValidateToken(token string) error {
	return nil
}

func (L *Logger) Log(level slog.Level,ctx context.Context, message string, args ...any) {
	event := Event{

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

func (L *Logger) Run() {
	batch := make([]Event, 0, 100)
	for {
		event := <-L.channel
		batch = append(batch, event)
		if len(batch) >= 100 {
			L.SendBatch(batch)
			batch = make([]Event, 0,100)	
		}
	}
}

func (L *Logger) SendBatch(batch []Event) {
	// Implement the logic to send the batch of events to the desired destination (e.g., a logging service or database).
}