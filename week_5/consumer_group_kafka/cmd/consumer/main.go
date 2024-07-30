package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

const (
	brokerAddress = "localhost:9092"
	groupID       = "1"
	topicName     = "test-topic"
)

func main() {
	keepRunning := true
	log.Println("starting a new Sarama consumer")

	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := Consumer{
		ready: make(chan bool),
	}

	client, err := sarama.NewConsumerGroup(strings.Split(brokerAddress, ","), groupID, config)
	if err != nil {
		log.Fatalf("failed to create consumer group client: %v\n", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		consume(ctx, client, consumer)
	}()

	<-consumer.ready
	log.Println("sarama consumer up and running ...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		log.Fatalf("failed to close consumer group client: %v", err.Error())
	}
}

func consume(ctx context.Context, client sarama.ConsumerGroup, consumer Consumer) {
	for {
		err := client.Consume(ctx, strings.Split(topicName, ","), &consumer)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			}

			log.Fatalf("failed to consume: %v", err.Error())
		}

		if ctx.Err() != nil {
			return
		}

		log.Printf("rebalancing\n")
		consumer.ready = make(chan bool)
	}
}
