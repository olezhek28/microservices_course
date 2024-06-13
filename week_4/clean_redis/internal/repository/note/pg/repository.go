package pg

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/olezhek28/platform_common/pkg/db"
	"github.com/pkg/errors"

	"github.com/olezhek28/microservices_course/week_4/clean_redis/internal/model"
	"github.com/olezhek28/microservices_course/week_4/clean_redis/internal/repository"
	"github.com/olezhek28/microservices_course/week_4/clean_redis/internal/repository/note/pg/converter"
	modelRepo "github.com/olezhek28/microservices_course/week_4/clean_redis/internal/repository/note/pg/model"
)

const (
	tableName = "note"

	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "content"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.NoteRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn, contentColumn).
		Values(info.Title, info.Content).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "note_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	builder := sq.Select(idColumn, titleColumn, contentColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "note_repository.Get",
		QueryRaw: query,
	}

	var note modelRepo.Note
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrorNoteNotFound
		}

		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
