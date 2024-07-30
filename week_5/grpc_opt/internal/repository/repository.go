package repository

import (
	"context"

	"github.com/olezhek28/microservices_course/week_5/grpc_opt/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}
