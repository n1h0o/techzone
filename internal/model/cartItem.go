package model

// CartItem представляет товар в корзине пользователя
type CartItem struct {
	ID        int64 `json:"id"`
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
