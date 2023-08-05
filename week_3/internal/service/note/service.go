package note

import (
	"github.com/olezhek28/microservices_course/week_3/internal/repository"
	"github.com/olezhek28/microservices_course/week_3/internal/service"
)

type serv struct {
	noteRepository repository.NoteRepository
}

func NewService(noteRepository repository.NoteRepository) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
	}
}
