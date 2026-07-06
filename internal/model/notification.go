package model

import "time"

// Notification представляет уведомление пользователя
type Notification struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	OrderID   int64     `json:"order_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
