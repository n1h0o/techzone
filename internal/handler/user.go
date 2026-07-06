package handler

import (
	"encoding/json"
	"net/http"
	"techzone/internal/middleware"
	"techzone/pkg/jwt"
)

// GetMe godoc
//
// @Summary Получить профиль
// @Description Возвращает информацию о текущем пользователе
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} handler.MeResponse
// @Failure 401 {string} string
// @Router /me [get]
func GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	resp := MeResponse{
		UserID: claims.UserID,
		Login:  claims.Login,
		Role:   claims.Role,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
