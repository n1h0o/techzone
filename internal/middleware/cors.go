package middleware

import (
	"net/http"
	"os"
)

func CORSMiddleware(next http.Handler) http.Handler {
	frontend := os.Getenv("FRONTEND_URL")

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if frontend != "" {
			w.Header().Set(
				"Access-Control-Allow-Origin",
				frontend,
			)
		}

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization, Idempotency-Key",
		)

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
