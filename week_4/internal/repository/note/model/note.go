package model

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64
	Info      NoteInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type NoteInfo struct {
	Title   string
	Content string
}
