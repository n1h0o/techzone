package handler

import (
	"encoding/json"
	"net/http"
	notificationClient "techzone/internal/grpcclient/notification"
	"techzone/internal/middleware"
	"techzone/internal/notification/pb"
	"techzone/pkg/jwt"
)

type NotificationHandler struct {
	notificationClient *notificationClient.Client
}

func NewNotificationHandler(
	notificationClient *notificationClient.Client,
) *NotificationHandler {
	return &NotificationHandler{
		notificationClient: notificationClient,
	}
}

// GetNotifications godoc
//
// @Summary Получить уведомления
// @Description Возвращает список уведомлений текущего пользователя
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Success 200 {object} handler.NotificationsResponse
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /notifications [get]
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

	resp, err := h.notificationClient.GetNotifications(
		r.Context(),
		&pb.GetNotificationsRequest{
			UserId: claims.UserID,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notifications := make([]NotificationResponse, 0, len(resp.Notifications))

	for _, n := range resp.Notifications {
		notifications = append(notifications, NotificationResponse{
			ID:        n.Id,
			UserID:    n.UserId,
			OrderID:   n.OrderId,
			Message:   n.Message,
			CreatedAt: n.CreatedAt.AsTime(),
		},
		)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		NotificationsResponse{
			Notifications: notifications,
		},
	); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
