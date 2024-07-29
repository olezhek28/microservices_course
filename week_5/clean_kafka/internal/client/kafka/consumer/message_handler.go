package consumer

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

type GroupHandler struct {
	msgHandler Handler
}

func NewGroupHandler() *GroupHandler {
	return &GroupHandler{}
}

// Setup запускается в начале новой сессии до вызова ConsumeClaim
func (c *GroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup запускается в конце жизни сессии после того как все горутины ConsumeClaim завершаться
func (c *GroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim должен запустить потребительский цикл сообщений ConsumerGroupClaim().
// После закрытия канала Messages() обработчик должен завершить обработку
func (c *GroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Код ниже не стоит перемещать в горутину, так как ConsumeClaim
	// уже запускается в горутине, см.:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L869
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed\n")
				return nil
			}

			log.Printf("message claimed: value = %s, timestamp = %v, topic = %s\n", string(message.Value), message.Timestamp, message.Topic)

			err := c.msgHandler(session.Context(), message)
			if err != nil {
				log.Printf("error handling message: %v\n", err)
				continue
			}

			session.MarkMessage(message, "")

		// Должен вернуться, когда `session.Context()` завершен.
		// В противном случае возникнет `ErrRebalanceInProgress` или `read tcp <ip>:<port>: i/o timeout` при перебалансировке кафки. см.:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			log.Printf("session context done\n")
			return nil
		}
	}
}
