package model

import "time"

// Order представляет заказ пользователя
type Order struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}
