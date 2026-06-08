package service

type AddToCartInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
