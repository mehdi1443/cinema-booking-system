package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type RedisLimiter struct {
	limiter   *redis_rate.Limiter
	appLogger domainlogger.Logger
}

func NewRedisLimiter(redisClient database.Cache, appLogger domainlogger.Logger) *RedisLimiter {
	return &RedisLimiter{
		limiter:   redis_rate.NewLimiter(redisClient.GetRDB()),
		appLogger: appLogger,
	}
}

func (rlm *RedisLimiter) PerSecond(limit int) gin.HandlerFunc {
	return rlm.buildMiddleware(redis_rate.PerSecond(limit))
}

func (rlm *RedisLimiter) PerMinute(limit int) gin.HandlerFunc {
	return rlm.buildMiddleware(redis_rate.PerMinute(limit))
}

func (rlm *RedisLimiter) buildMiddleware(limitRule redis_rate.Limit) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIdentifier := c.Request.RemoteAddr
		if host, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
			clientIdentifier = host
		}
		if userID, exists := c.Get("user_id"); exists {
			clientIdentifier = fmt.Sprintf("user:%v", userID)
		}

		key := fmt.Sprintf("rate_limit:%s:%s", c.FullPath(), clientIdentifier)

		res, err := rlm.limiter.Allow(c.Request.Context(), key, limitRule)
		if err != nil {
			c.Error(err)
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(limitRule.Rate))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed == 0 {
			retryAfterSec := int(res.RetryAfter.Seconds())
			if retryAfterSec < 1 {
				retryAfterSec = 1
			}
			c.Header("Retry-After", strconv.Itoa(retryAfterSec))

			if rlm.appLogger != nil {
				rlm.appLogger.Warn(c.Request.Context(), "rate limit exceeded", domainlogger.String("path", c.FullPath()))
			}

			response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later...")
			c.Abort()
			return
		}

		c.Next()
	}
}
