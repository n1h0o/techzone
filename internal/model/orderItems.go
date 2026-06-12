package model

type OrderItem struct {
	ID        int64
	OrderID   int64
	ProductID int64
	Quantity  int
	Price     float64
}
