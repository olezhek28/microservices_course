package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/olezhek28/microservices_course/week_2/config/internal/repository/note/model"
	desc "github.com/olezhek28/microservices_course/week_2/config/pkg/note_v1"
)

func ToNoteFromRepo(note *model.Note) *desc.Note {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.Note{
		Id:        note.ID,
		Info:      ToNoteInfoFromRepo(note.Info),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToNoteInfoFromRepo(info *model.Info) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
