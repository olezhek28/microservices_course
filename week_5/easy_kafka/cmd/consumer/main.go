package main

import (
	"log"

	"github.com/IBM/sarama"
)

const (
	brokerAddress = "localhost:9092"
	topicName     = "test-topic"
)

func main() {
	consumer, err := newConsumer([]string{brokerAddress})
	if err != nil {
		log.Fatalf("failed to start consumer: %v\n", err.Error())
	}

	defer func() {
		if err = consumer.Close(); err != nil {
			log.Fatalf("Failed to close consumer: %v\n", err.Error())
		}
	}()

	consumeMessages(consumer, topicName)
}

func newConsumer(brokerList []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func consumeMessages(consumer sarama.Consumer, topic string) {
	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to start consumer for partition %d: %v", 0, err)
	}

	defer func() {
		err = pc.Close()
		if err != nil {
			log.Fatalf("failed to close partition consumer: %v\n", err.Error())
		}
	}()

	for {
		select {
		// Чтение сообщения из Kafka
		case msg, ok := <-pc.Messages():
			if !ok {
				log.Println("channel closed, exiting goroutine")
				return
			}

			log.Printf("message received: %v\n", string(msg.Value))
		}
	}
}
