package event

type PaymentCompletedEvent struct {
	PaymentID int64 `json:"payment_id"`
	OrderID   int64 `json:"order_id"`
	UserID    int64 `json:"user_id"`
}
