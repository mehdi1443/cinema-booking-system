package logger

import "context"

type contextKey string

const (
	requestIDContextKey contextKey = "request_id"
	traceIDContextKey   contextKey = "trace_id"
)

func ContextWithRequestID(
	ctx context.Context,
	requestID string,
) context.Context {
	return context.WithValue(
		ctx,
		requestIDContextKey,
		requestID,
	)
}

func RequestIDFromContext(
	ctx context.Context,
) (string, bool) {
	value := ctx.Value(requestIDContextKey)

	if value == nil {
		return "", false
	}

	requestID, ok := value.(string)

	return requestID, ok && requestID != ""
}

func ContextWithTraceID(
	ctx context.Context,
	traceID string,
) context.Context {
	return context.WithValue(
		ctx,
		traceIDContextKey,
		traceID,
	)
}

func TraceIDFromContext(
	ctx context.Context,
) (string, bool) {
	value := ctx.Value(traceIDContextKey)

	if value == nil {
		return "", false
	}

	traceID, ok := value.(string)

	return traceID, ok && traceID != ""
}
