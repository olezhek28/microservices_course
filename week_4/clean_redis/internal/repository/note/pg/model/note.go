package model

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
