package repository

import (
	"context"

	desc "github.com/olezhek28/microservices_course/week_3/pkg/note_v1"
)

type NoteRepository interface {
	Create(ctx context.Context, info *desc.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*desc.Note, error)
}
