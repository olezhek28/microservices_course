package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/olezhek28/microservices_course/week_4/internal/model"
	"github.com/olezhek28/microservices_course/week_4/internal/repository"
	repoMocks "github.com/olezhek28/microservices_course/week_4/internal/repository/mocks"
	"github.com/olezhek28/microservices_course/week_4/internal/service/note"
)

func TestGet(t *testing.T) {
	type noteRepositoryMockFunc func(mc *minimock.Controller) repository.NoteRepository

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		title     = gofakeit.Animal()
		content   = gofakeit.Animal()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		repoErr = fmt.Errorf("repo error")

		res = &model.Note{
			ID: id,
			Info: model.NoteInfo{
				Title:   title,
				Content: content,
			},
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.Note
		err                error
		noteRepositoryMock noteRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			noteRepositoryMock: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repoMocks.NewNoteRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  repoErr,
			noteRepositoryMock: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repoMocks.NewNoteRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			noteRepoMock := tt.noteRepositoryMock(mc)
			service := note.NewMockService(noteRepoMock)

			newID, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
