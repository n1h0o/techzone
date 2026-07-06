package kafka

import (
	"context"
	"encoding/json"
	"log"
	"techzone/internal/event"
	"techzone/internal/service"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client

	workerPool *service.NotificationWorkerPool
}

func NewConsumer(
	client *kgo.Client,
	workerPool *service.NotificationWorkerPool,
) *Consumer {
	return &Consumer{
		client:     client,
		workerPool: workerPool,
	}
}

func (c *Consumer) Start(
	ctx context.Context,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		fetches := c.client.PollFetches(ctx)

		if ctx.Err() != nil {
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("kafka error: %v", err)
			}
			continue
		}
		fetches.EachRecord(func(record *kgo.Record) {
			if err := c.handleRecord(ctx, record); err != nil {
				log.Printf("failed to handle record %v", err)
			}
		})
	}
}

func (c *Consumer) handleOrderCreated(
	ctx context.Context,
	record *kgo.Record,
) error {
	var evt event.OrderCreatedEvent

	if err := json.Unmarshal(
		record.Value,
		&evt,
	); err != nil {
		log.Printf("failed to decode event: %v", err)
		return err
	}

	c.workerPool.Submit(
		service.NotificationJob{
			UserID:  evt.UserID,
			OrderID: evt.OrderID,
			Message: "Заказ успешно создан",
		},
	)

	return c.client.CommitRecords(ctx, record)
}

func (c *Consumer) handlePaymentCreated(
	ctx context.Context,
	record *kgo.Record,
) error {
	var evt event.PaymentCompletedEvent

	if err := json.Unmarshal(
		record.Value,
		&evt,
	); err != nil {
		log.Printf("failed to decode event: %v", err)
		return err
	}

	c.workerPool.Submit(
		service.NotificationJob{
			UserID:  evt.UserID,
			OrderID: evt.OrderID,
			Message: "Оплата заказа успешно выполнена",
		},
	)

	return c.client.CommitRecords(ctx, record)
}

func (c *Consumer) handleRecord(
	ctx context.Context,
	record *kgo.Record,
) error {

	switch record.Topic {
	case "order.created":
		return c.handleOrderCreated(ctx, record)

	case "payment.completed":
		return c.handlePaymentCreated(ctx, record)

	default:
		log.Printf("unknown topic: %s", record.Topic)
		return nil
	}

}
