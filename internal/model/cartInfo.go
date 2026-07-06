package model

// CartItemInfo представляет информацию о товаре в корзине пользователя
type CartItemInfo struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}
