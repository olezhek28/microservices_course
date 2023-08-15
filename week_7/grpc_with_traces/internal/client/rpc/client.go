package rpc

import (
	"context"

	"github.com/olezhek28/microservices_course/week_7/grpc_with_traces/internal/model"
)

type OtherServiceClient interface {
	Get(ctx context.Context, id int64) (*model.Note, error)
}
