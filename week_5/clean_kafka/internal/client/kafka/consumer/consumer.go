package consumer

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type consumer struct {
	topicName            string
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *GroupHandler
}

func NewConsumer(
	consumerGroup sarama.ConsumerGroup,
	consumerGroupHandler *GroupHandler,
) *consumer {
	return &consumer{
		consumerGroup:        consumerGroup,
		consumerGroupHandler: consumerGroupHandler,
	}
}

func (c *consumer) Consume(ctx context.Context, topicName string, handler Handler) (err error) {
	c.topicName = topicName
	c.consumerGroupHandler.msgHandler = handler

	return c.consume(ctx)
}

func (c *consumer) Close() error {
	return c.consumerGroup.Close()
}

func (c *consumer) consume(ctx context.Context) error {
	for {
		err := c.consumerGroup.Consume(ctx, strings.Split(c.topicName, ","), c.consumerGroupHandler)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		log.Printf("rebalancing...\n")
	}
}
