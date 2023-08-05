package model

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64        `db:"id"`
	Info      *Info        `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type Info struct {
	Title   string `db:"title"`
	Content string `db:"content"`
}
