package middleware

import (
	"net/http"
	"techzone/pkg/jwt"
)

// пропускает дальше только администраторов
func AdminMiddleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			claims, ok := r.Context().Value(UserKey).(*jwt.Claims)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if claims.Role != "admin" {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}
