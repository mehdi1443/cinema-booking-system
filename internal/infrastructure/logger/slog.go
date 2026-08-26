package logger

import (
	"context"
	"log/slog"
	"os"
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

	handler := slog.NewJSONHandler(
		os.Stdout,
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

func (l *SlogLogger) Debug(
	ctx context.Context,
	msg string,
	fields ...logger.Field,
) {
	l.log(ctx, slog.LevelDebug, msg, fields...)
}

func (l *SlogLogger) Info(
	ctx context.Context,
	msg string,
	fields ...logger.Field,
) {
	l.log(ctx, slog.LevelInfo, msg, fields...)
}

func (l *SlogLogger) Warn(
	ctx context.Context,
	msg string,
	fields ...logger.Field,
) {
	l.log(ctx, slog.LevelWarn, msg, fields...)
}

func (l *SlogLogger) Error(
	ctx context.Context,
	msg string,
	fields ...logger.Field,
) {
	l.log(ctx, slog.LevelError, msg, fields...)
}

func (l *SlogLogger) With(fields ...logger.Field) logger.Logger {
	return &SlogLogger{
		logger: l.logger.With(toSlogAttrs(fields)...),
	}
}

func (l *SlogLogger) log(
	ctx context.Context,
	level slog.Level,
	msg string,
	fields ...logger.Field,
) {
	args := make([]any, 0, len(fields)*2+4)

	for _, field := range fields {
		if field.Key == "" {
			continue
		}

		args = append(args, field.Key, sanitizeValue(field.Key, field.Value))
	}

	// Request correlation.
	if requestID, ok := RequestIDFromContext(ctx); ok {
		args = append(args, "request_id", requestID)
	}

	// Trace correlation.
	if traceID, ok := TraceIDFromContext(ctx); ok {
		args = append(args, "trace_id", traceID)
	}

	l.logger.Log(
		ctx,
		level,
		msg,
		args...,
	)
}

func toSlogAttrs(fields []logger.Field) []any {
	args := make([]any, 0, len(fields)*2)

	for _, field := range fields {
		if field.Key == "" {
			continue
		}

		args = append(
			args,
			field.Key,
			sanitizeValue(field.Key, field.Value),
		)
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
		"password",
		"passwd",
		"secret",
		"token",
		"access_token",
		"refresh_token",
		"authorization",
		"cookie",
		"set-cookie",
		"api_key",
		"apikey",
		"private_key",
	}

	for _, sensitive := range sensitiveKeys {
		if key == sensitive {
			return true
		}
	}

	return false
}
