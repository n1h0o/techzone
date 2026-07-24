package kafka

import (
	"context"
	"encoding/json"
	"techzone/internal/event"
	"techzone/internal/metrics"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
}

// оборачивает kafka клиент в интерфейс доменного publisher слоя
func NewProducer(
	client *kgo.Client,
) *Producer {
	return &Producer{
		client: client,
	}
}

// публикует событие создания заказа после успешного коммита
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

	err = p.client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return err
	}

	metrics.KafkaMessagesProducedTotal.Inc()
	return nil
}

// публикует событие успешной оплаты
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
