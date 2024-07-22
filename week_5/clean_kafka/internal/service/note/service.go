package note

import (
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/service"
)

type serv struct {
}

func NewService() service.NoteService {
	return &serv{}
}
