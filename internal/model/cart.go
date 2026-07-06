package model

import "time"

// Cart представляет корзину пользователя
type Cart struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
