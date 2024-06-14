package converter

import (
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/model"
	modelRepo "github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository/note/pg/model"
)

func ToNoteFromRepo(note *modelRepo.Note) *model.Note {
	return &model.Note{
		ID: note.ID,
		Info: model.NoteInfo{
			Title:   note.Title,
			Content: note.Content,
		},
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
