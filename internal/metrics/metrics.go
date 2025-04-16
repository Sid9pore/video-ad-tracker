package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define metrics
var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func InitMetrics() {
	// Register metrics with Prometheus
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestDuration)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Record metrics
		totalRequests.WithLabelValues(r.Method, r.RequestURI).Inc()
		requestDuration.WithLabelValues(r.Method, r.RequestURI).Observe(time.Since(start).Seconds())
	})
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
