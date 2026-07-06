package model

// OrderItemInfo содержит информацию о товаре в заказе
type OrderItemInfo struct {
	ProductID int64   `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}
