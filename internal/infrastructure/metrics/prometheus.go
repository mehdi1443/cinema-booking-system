package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusMetrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cinema_http_requests_total",
				Help: "Total number of HTTP requests processed by the cinema booking API.",
			},
			[]string{"method", "route", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "cinema_http_request_duration_seconds",
				Help:    "Histogram of response latency for HTTP requests in seconds.",
				Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
			},
			[]string{"method", "route"},
		),
	}
}

func (pm *PrometheusMetrics) IncHTTPRequest(method, route, status string) {
	pm.httpRequestsTotal.WithLabelValues(method, route, status).Inc()
}

func (pm *PrometheusMetrics) ObserveHTTPRequestDuration(method, route string, duration float64) {
	pm.httpRequestDuration.WithLabelValues(method, route).Observe(duration)
}
