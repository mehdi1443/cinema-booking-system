package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMiddleware struct {
	metricsClient metrics.MetricsClient
}

func NewPrometheusMiddleware(metricsClient metrics.MetricsClient) *PrometheusMiddleware {
	return &PrometheusMiddleware{
		metricsClient: metricsClient,
	}
}

func (pm *PrometheusMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// استفاده از FullPath به جای URL.Path برای جلوگیری از انفجار حافظه پرومتیوس
		route := c.FullPath()
		if route == "" {
			route = "unknown_route"
		}

		method := c.Request.Method

		pm.metricsClient.IncHTTPRequest(method, route, status)
		pm.metricsClient.ObserveHTTPRequestDuration(method, route, duration)
	}
}

func SetupPrometheusEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
