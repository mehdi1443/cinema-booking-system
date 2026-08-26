package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/logger"
)

// GinLoggerMiddleware injects RequestID/TraceID and logs the HTTP request details.
func GinLoggerMiddleware(log domainlogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := c.Request.Context()
		ctx = logger.ContextWithTraceID(ctx, traceID)
		ctx = logger.ContextWithRequestID(ctx, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		fields := []domainlogger.Field{
			domainlogger.String("method", c.Request.Method),
			domainlogger.String("path", c.Request.URL.Path),
			domainlogger.String("ip", c.ClientIP()),
			domainlogger.Int("status", statusCode),
			domainlogger.String("latency", latency.String()),
			domainlogger.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, domainlogger.String("gin_errors", c.Errors.String()))
		}

		msg := "HTTP Request completed"

		if statusCode >= 500 {
			log.Error(ctx, msg, fields...)
		} else if statusCode >= 400 {
			log.Warn(ctx, msg, fields...)
		} else {
			log.Info(ctx, msg, fields...)
		}
	}
}
