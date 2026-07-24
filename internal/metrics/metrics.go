package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	OrdersCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "total number of created orders",
		},
	)

	NotificationsCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notifications_created_total",
			Help: "Total number of created notifications",
		},
	)

	KafkaMessagesProducedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_produced_total",
			Help: "Total produced kafka messages",
		},
	)

	KafkaMessagesConsumedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total produced kafka messages",
		},
	)
)

var HTTPRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP request latency",
		Buckets: []float64{
			0.005,
			0.01,
			0.025,
			0.05,
			0.1,
			0.25,
			0.5,
			1,
			2.5,
			5,
		},
	},
	[]string{"method", "path", "status"},
)

func Init() {
	prometheus.MustRegister(
		OrdersCreatedTotal,
		NotificationsCreatedTotal,
		KafkaMessagesProducedTotal,
		KafkaMessagesConsumedTotal,
		HTTPRequestDuration,
	)
}
