package kafka

import (
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

// создает клиент только для публикации событий
func NewProducerClient() (*kgo.Client, error) {
	return kgo.NewClient(
		kgo.SeedBrokers(
			os.Getenv("KAFKA_BROKERS"),
		),
	)
}

// создает consumer клиента для notification сервиса
func NewConsumerClient() (*kgo.Client, error) {
	return kgo.NewClient(
		kgo.SeedBrokers(
			os.Getenv("KAFKA_BROKERS"),
		),
		kgo.ConsumeTopics(
			"order.created",
			"payment.completed",
		),
		kgo.ConsumerGroup(
			"techzone-notifications",
		),
	)
}
