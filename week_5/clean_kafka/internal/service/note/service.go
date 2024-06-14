package note

import (
	"github.com/olezhek28/platform_common/pkg/db"

	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/service"
)

type serv struct {
	noteRepository repository.NoteRepository
	txManager      db.TxManager
}

func NewService(
	noteRepository repository.NoteRepository,
	txManager db.TxManager,
) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
		txManager:      txManager,
	}
}
