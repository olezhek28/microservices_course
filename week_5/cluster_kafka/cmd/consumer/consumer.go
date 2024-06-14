package main

import (
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	ready chan bool
}

// Setup запускается в начале новой сессии до вызова ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Помечаем консьюмер как готовый к работе
	close(c.ready)
	return nil
}

// Cleanup запускается в конце жизни сессии is run at the end of a session после того как все горутины ConsumeClaim завершаться
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim должен запустить потребительский цикл сообщений ConsumerGroupClaim().
// После закрытия канала Messages() обработчик должен завершить обработку
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Код ниже не стоит перемещать в горутину, так как ConsumeClaim
	// уже запускается в горутине, см.:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed\n")
				return nil
			}

			log.Printf("message claimed: value = %s, timestamp = %v, topic = %s\n", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")

		// Должен вернуться, когда `session.Context()` завершен.
		// В противном случае возникнет `ErrRebalanceInProgress` или `read tcp <ip>:<port>: i/o timeout` при перебалансировке кафки. см.:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
