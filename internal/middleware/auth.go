package middleware

import (
	"context"
	"net/http"
	"strings"
	"techzone/internal/config"
	"techzone/pkg/jwt"
)

// проверяет bearer токен и кладет claims в контекст запроса
func AuthMiddleware(
	cfg *config.Config,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")

				if authHeader == "" {
					http.Error(w, "missing authorization header", http.StatusUnauthorized)
					return
				}
				parts := strings.Split(authHeader, " ")

				if len(parts) != 2 || parts[0] != "Bearer" {
					http.Error(w, "invalid authorization header", http.StatusUnauthorized)
					return
				}
				tokenString := parts[1]

				// разбирает токен один раз чтобы downstream код работал только с claims
				claims, err := jwt.ParseToken(
					tokenString,
					cfg.JWTSecret,
				)
				if err != nil {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(
					r.Context(),
					UserKey,
					claims,
				)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
			},
		)
	}
}
