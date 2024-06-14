package pg

import (
	"context"
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/client/cache"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/model"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository/note/redis/converter"
	modelRepo "github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/repository/note/redis/model"
)

type repo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) repository.NoteRepository {
	return &repo{cl: cl}
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	id := int64(1)

	note := modelRepo.Note{
		ID:          id,
		Title:       info.Title,
		Content:     info.Content,
		CreatedAtNs: time.Now().UnixNano(),
	}

	idStr := strconv.FormatInt(id, 10)
	err := r.cl.HashSet(ctx, idStr, note)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorNoteNotFound
	}

	var note modelRepo.Note
	err = redigo.ScanStruct(values, &note)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
