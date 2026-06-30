package event

type OrderCreatedEvent struct {
	OrderID int64 `json:"order_id"`
	UserID  int64 `json:"user_id"`
}
