package logger

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
)

const (
	requestIDKey = "request_id"
	traceIDKey   = "trace_id"
)

type SlogLogger struct {
	logger *slog.Logger
}

func New(cfg Config) logger.Logger {
	handlerOptions := &slog.HandlerOptions{
		Level:     cfg.slogLevel(),
		AddSource: cfg.AddSource,
	}

	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("failed to create log directory: %v", err)
	}

	logFilePath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	var writer io.Writer = os.Stdout
	if err != nil {
		log.Printf("failed to open log file, falling back to stdout only: %v", err)
	} else {
		writer = io.MultiWriter(os.Stdout, logFile)
	}

	handler := slog.NewJSONHandler(
		writer,
		handlerOptions,
	)

	baseLogger := slog.New(handler).With(
		slog.String("service", cfg.ServiceName),
		slog.String("environment", cfg.Environment),
	)

	return &SlogLogger{
		logger: baseLogger,
	}
}

func (l *SlogLogger) Debug(ctx context.Context, msg string, fields ...logger.Field) {
	l.log(ctx, slog.LevelDebug, msg, fields...)
}

func (l *SlogLogger) Info(ctx context.Context, msg string, fields ...logger.Field) {
	l.log(ctx, slog.LevelInfo, msg, fields...)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string, fields ...logger.Field) {
	l.log(ctx, slog.LevelWarn, msg, fields...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string, fields ...logger.Field) {
	l.log(ctx, slog.LevelError, msg, fields...)
}

func (l *SlogLogger) With(fields ...logger.Field) logger.Logger {
	return &SlogLogger{
		logger: l.logger.With(toSlogAttrs(fields)...),
	}
}

func (l *SlogLogger) log(ctx context.Context, level slog.Level, msg string, fields ...logger.Field) {
	args := make([]any, 0, len(fields)*2+4)

	for _, field := range fields {
		if field.Key == "" {
			continue
		}
		args = append(args, field.Key, sanitizeValue(field.Key, field.Value))
	}

	if requestID, ok := RequestIDFromContext(ctx); ok {
		args = append(args, requestIDKey, requestID)
	}

	if traceID, ok := TraceIDFromContext(ctx); ok {
		args = append(args, traceIDKey, traceID)
	}

	l.logger.Log(ctx, level, msg, args...)
}

func toSlogAttrs(fields []logger.Field) []any {
	args := make([]any, 0, len(fields)*2)

	for _, field := range fields {
		if field.Key == "" {
			continue
		}
		args = append(args, field.Key, sanitizeValue(field.Key, field.Value))
	}

	return args
}

func sanitizeValue(key string, value any) any {
	if isSensitiveKey(key) {
		return "[REDACTED]"
	}
	return value
}

func isSensitiveKey(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))

	sensitiveKeys := []string{
		"password", "passwd", "secret", "token",
		"access_token", "refresh_token", "authorization",
		"cookie", "set-cookie", "api_key", "apikey", "private_key",
	}

	for _, sensitive := range sensitiveKeys {
		if key == sensitive || strings.HasSuffix(key, "_"+sensitive) {
			return true
		}
	}
	// Cover common application-specific variants such as otp_code,
	// password_hash, and idempotency_key without maintaining an unbounded list.
	for _, suffix := range []string{"_otp", "_code", "_password", "_secret", "_token", "_key"} {
		if strings.HasSuffix(key, suffix) {
			return true
		}
	}

	return false
}
