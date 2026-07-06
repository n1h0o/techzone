package payment

import "context"

type PaymentRequest struct {
	OrderID int64
	UserID  int64
	Amount  int64
}

type PaymentResponse struct {
	TransactionID string
	Status        string
}

type Gateway interface {
	Pay(
		ctx context.Context,
		req PaymentRequest,
	) (*PaymentResponse, error)
}
