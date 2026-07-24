package middleware

import (
	"net/http"
	"strconv"
	"techzone/internal/metrics"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pattern := r.Pattern
		if pattern == "" {
			pattern = "unknown"
		}

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		start := time.Now()

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		metrics.HTTPRequestDuration.
			WithLabelValues(
				r.Method,
				pattern,
				strconv.Itoa(rw.status),
			).
			Observe(duration)
	})
}
