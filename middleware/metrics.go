// Port https://github.com/zbindenren/negroni-prometheus for chi router
package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	defBuckets = []float64{100, 200, 500, 1000}
)

const (
	reqsName    = "chi_requests_total"
	latencyName = "chi_request_duration_milliseconds"
)

// MetricsMiddleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type MetricsMiddleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

func SetBuckets(b []float64) {
	defBuckets = b
}

// NewPrometheus returns a new prometheus MetricsMiddleware handler.
func NewPrometheus(name string, buckets ...float64) func(next http.Handler) http.Handler {
	var m MetricsMiddleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = defBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)
	return m.handler
}

func (c MetricsMiddleware) handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		c.reqs.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.Path).Inc()
		c.latency.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}
