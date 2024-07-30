package kafka

import (
	"context"

	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/client/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
