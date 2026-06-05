package handler

import (
	"encoding/json"
	"net/http"
	"techzone/internal/middleware"
	"techzone/pkg/jwt"
)

type MeResponse struct {
	UserID int64  `json:"user_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}

func GetMe(w http.ResponseWriter, r *http.Request) {
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
