package note

import (
	"context"
	"log"

	"github.com/olezhek28/microservices_course/week_5/internal/converter"
	desc "github.com/olezhek28/microservices_course/week_5/pkg/note_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	noteObj, err := i.noteService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", noteObj.ID, noteObj.Info.Title, noteObj.Info.Content, noteObj.CreatedAt, noteObj.UpdatedAt)

	return &desc.GetResponse{
		Note: converter.ToNoteFromService(noteObj),
	}, nil
}
