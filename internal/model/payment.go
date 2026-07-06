package model

import "time"

// Payment представляет информацию об оплате заказа
type Payment struct {
	ID             int64     `json:"id"`
	OrderID        int64     `json:"order_id"`
	UserID         int64     `json:"user_id"`
	Amount         int64     `json:"amount"`
	Status         string    `json:"status"`
	IdempotencyKey string    `json:"idempotency_key"`
	TransactionID  string    `json:"transaction_id"`
	CreatedAt      time.Time `json:"created_at"`
}
