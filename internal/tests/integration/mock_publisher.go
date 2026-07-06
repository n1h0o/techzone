package integration

import (
	"context"
	"techzone/internal/event"
)

type MockPublisher struct{}

func (m *MockPublisher) PublishOrderCreated(
	ctx context.Context,
	event event.OrderCreatedEvent,
) error {
	return nil
}

func (m *MockPublisher) PublishPaymentCompleted(
	ctx context.Context,
	event event.PaymentCompletedEvent,
) error {
	return nil
}
