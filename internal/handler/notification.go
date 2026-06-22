package handler

import (
	"encoding/json"
	"net/http"
	"techzone/internal/middleware"
	"techzone/internal/service"
	"techzone/pkg/jwt"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(
	notificationService *service.NotificationService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) GetNotifications(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*jwt.Claims)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := h.notificationService.GetNotifications(
		r.Context(),
		claims.UserID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		map[string]any{
			"notifications": notifications,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
