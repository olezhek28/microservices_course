package interceptor

import (
	"context"

	"google.golang.org/grpc"

	"github.com/olezhek28/microservices_course/week_7/grpc_with_metrics/internal/metric"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	return handler(ctx, req)
}
