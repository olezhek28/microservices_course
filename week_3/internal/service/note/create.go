package note

import (
	"context"

	"github.com/olezhek28/microservices_course/week_3/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	return 0, nil
}
