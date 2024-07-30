package note_saver

import (
	"context"

	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/client/kafka"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository"
	def "github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/service"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	noteRepository repository.NoteRepository
	consumer       kafka.Consumer
}

func NewService(
	noteRepository repository.NoteRepository,
	consumer kafka.Consumer,
) *service {
	return &service{
		noteRepository: noteRepository,
		consumer:       consumer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, "test-topic", s.NoteSaveHandler)
	}()

	return errChan
}
