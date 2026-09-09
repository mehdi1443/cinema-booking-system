package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

func Audit(repo DB.OperationsRepositoryInterface, log domainlogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if repo == nil {
			c.Next()
			return
		}
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		requestID := truncateAuditValue(strings.TrimSpace(c.GetHeader("X-Request-ID")), 100)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Header("X-Request-ID", requestID)
		started := time.Now()
		c.Next()

		var actorID *int64
		if value, exists := c.Get(domainjwt.CtxUserIDKey); exists {
			if id, ok := value.(int64); ok {
				actorID = &id
			}
		}
		routeValue := c.FullPath()
		var route *string
		if routeValue != "" {
			route = &routeValue
		}
		pathValue := truncateAuditValue(c.Request.URL.Path, 500)
		if route != nil {
			value := truncateAuditValue(*route, 500)
			route = &value
		}
		ip, userAgent, rid := truncateAuditValue(c.ClientIP(), 64), truncateAuditValue(c.Request.UserAgent(), 500), requestID
		entry := &entity.AuditLog{ActorUserID: actorID, Method: c.Request.Method, Path: pathValue, Route: route, StatusCode: c.Writer.Status(), IPAddress: &ip, UserAgent: &userAgent, RequestID: &rid, DurationMS: time.Since(started).Milliseconds()}
		auditCtx, cancel := context.WithTimeout(context.WithoutCancel(c.Request.Context()), 2*time.Second)
		defer cancel()
		if err := repo.RecordAudit(auditCtx, entry); err != nil && log != nil {
			log.Error(auditCtx, "failed to persist audit log", domainlogger.Error(err), domainlogger.String("method", entry.Method), domainlogger.String("path", entry.Path))
		}
	}
}

func truncateAuditValue(value string, max int) string {
	runes := []rune(value)
	if len(runes) > max {
		return string(runes[:max])
	}
	return value
}
