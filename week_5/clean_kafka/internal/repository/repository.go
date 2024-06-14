package repository

import (
	"context"

	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}
