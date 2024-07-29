package model

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64        `json:"id"`
	Info      NoteInfo     `json:"info"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type NoteInfo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
