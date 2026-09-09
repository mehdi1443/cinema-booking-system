package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func Recovery(log domainlogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				if log != nil {
					log.Error(c.Request.Context(), "panic recovered from HTTP request",
						domainlogger.String("panic", fmt.Sprint(recovered)),
						domainlogger.String("stack", string(debug.Stack())),
					)
				}
				if !c.Writer.Written() {
					response.ErrorWithCode(c, http.StatusInternalServerError, string(exception.CodeInternal), "An unexpected internal error occurred")
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
