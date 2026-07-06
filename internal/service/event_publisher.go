package service

import (
	"context"
	"techzone/internal/event"
)

type EventPublisher interface {
	PublishOrderCreated(
		ctx context.Context,
		event event.OrderCreatedEvent,
	) error
	PublishPaymentCompleted(
		ctx context.Context,
		event event.PaymentCompletedEvent,
	) error
}
