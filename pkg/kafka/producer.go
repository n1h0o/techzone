package kafka

import (
	"context"
	"encoding/json"
	"techzone/internal/event"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
}

func NewProducer(
	client *kgo.Client,
) *Producer {
	return &Producer{
		client: client,
	}
}

func (p *Producer) PublishOrderCreated(
	ctx context.Context,
	event event.OrderCreatedEvent,
) error {

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Topic: "order.created",
		Value: data,
	}

	return p.client.ProduceSync(
		ctx,
		record,
	).FirstErr()
}

func (p *Producer) PublishPaymentCompleted(
	ctx context.Context,
	event event.PaymentCompletedEvent,
) error {

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Topic: "payment.completed",
		Value: data,
	}

	return p.client.ProduceSync(
		ctx,
		record,
	).FirstErr()
}
