package model

type Note struct {
	ID          int64  `redis:"id"`
	Title       string `redis:"title"`
	Content     string `redis:"content"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at"`
}
