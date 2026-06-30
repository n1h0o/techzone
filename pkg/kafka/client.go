package kafka

import (
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

func NewProducerClient() (*kgo.Client, error) {
	return kgo.NewClient(
		kgo.SeedBrokers(
			os.Getenv("KAFKA_BROKERS"),
		),
	)
}

func NewConsumerClient() (*kgo.Client, error) {
	return kgo.NewClient(
		kgo.SeedBrokers(
			os.Getenv("KAFKA_BROKERS"),
		),
		kgo.ConsumeTopics(
			"order.created",
		),
		kgo.ConsumerGroup(
			"techzone-notificatios",
		),
	)
}
