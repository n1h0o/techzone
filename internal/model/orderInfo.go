package model

type OrderInfo struct {
	ID         int64   `json:"id"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"total_price"`
}
