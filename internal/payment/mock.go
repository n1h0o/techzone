package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MockGateway struct{}

func NewMockGateway() *MockGateway {
	return &MockGateway{}
}

func (g *MockGateway) Pay(
	ctx context.Context,
	req PaymentRequest,
) (*PaymentResponse, error) {

	time.Sleep(500 * time.Millisecond)

	return &PaymentResponse{
		TransactionID: uuid.NewString(),
		Status:        StatusSuccess,
	}, nil
}
