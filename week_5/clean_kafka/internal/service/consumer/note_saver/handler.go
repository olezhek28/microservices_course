package note_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"

	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/model"
)

func (s *service) NoteSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	noteInfo := &model.NoteInfo{}
	err := json.Unmarshal(msg.Value, noteInfo)
	if err != nil {
		return err
	}

	id, err := s.noteRepository.Create(ctx, noteInfo)
	if err != nil {
		return err
	}

	log.Printf("Note with id %d created\n", id)

	return nil
}
