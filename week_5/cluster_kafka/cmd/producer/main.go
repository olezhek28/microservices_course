package main

import (
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/brianvoe/gofakeit/v6"
)

const (
	brokerAddress = "localhost:9092, localhost:9093, localhost:9094"
	topicName     = "test-topic"
)

func main() {
	producer, err := newSyncProducer(strings.Split(brokerAddress, ","))
	if err != nil {
		log.Fatalf("failed to start producer: %v\n", err.Error())
	}

	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("failed to close producer: %v\n", err.Error())
		}
	}()

	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(gofakeit.StreetName()),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message in Kafka: %v\n", err.Error())
		return
	}

	log.Printf("message sent to partition %d with offset %d\n", partition, offset)
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
