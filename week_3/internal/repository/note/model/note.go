package model

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64
	Info      Info
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type Info struct {
	Title   string
	Content string
}
