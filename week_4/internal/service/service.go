package service

import (
	"context"

	"github.com/olezhek28/microservices_course/week_4/internal/model"
)

type NoteService interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}
