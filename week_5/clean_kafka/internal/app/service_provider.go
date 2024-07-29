package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/olezhek28/platform_common/pkg/closer"
	"github.com/olezhek28/platform_common/pkg/db"
	"github.com/olezhek28/platform_common/pkg/db/pg"

	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/client/kafka"
	kafkaConsumer "github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/client/kafka/consumer"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/config"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/config/env"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository"
	noteRepo "github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository/note"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/service"
	noteSaverConsumer "github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/service/consumer/note_saver"
)

type serviceProvider struct {
	pgConfig            config.PGConfig
	kafkaConsumerConfig config.KafkaConsumerConfig

	dbClient db.Client

	noteRepository repository.NoteRepository

	noteSaverConsumer service.ConsumerService

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.ConsumerGroupHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := env.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
	if s.noteRepository == nil {
		s.noteRepository = noteRepo.NewRepository(s.DBClient(ctx))
	}

	return s.noteRepository
}

func (s *serviceProvider) NoteSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.noteSaverConsumer == nil {
		s.noteSaverConsumer = noteSaverConsumer.NewService(
			s.NoteRepository(ctx),
			s.Consumer(),
		)
	}

	return s.noteSaverConsumer
}

func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.ConsumerGroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewConsumerGroupHandler()
	}

	return s.consumerGroupHandler
}
