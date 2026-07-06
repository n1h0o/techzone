package model

// OrderInfo содержит краткую информацию о заказе
type OrderInfo struct {
	ID            int64   `json:"id"`
	Status        string  `json:"status"`
	PaymentStatus string  `json:"payment_status"`
	TotalPrice    float64 `json:"total_price"`
}
