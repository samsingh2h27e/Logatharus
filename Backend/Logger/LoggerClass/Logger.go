package Logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	terminal *slog.Logger
	ServiceName string
	isValid     bool
}

func NewLogger(serviceName string) *Logger {
	logger := &Logger{
		terminal : slog.New(slog.NewJSONHandler(os.Stdout,nil)),
		ServiceName: serviceName,
		isValid:     false,
	}
	return logger
}

func (L *Logger) ValidateToken(token string) error {
	return nil
}

func (L *Logger) Log(level slog.Level,ctx context.Context, message string, args ...any) {
	L.terminal.Log(ctx,level,message,args)

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