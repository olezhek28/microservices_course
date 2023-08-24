package interceptor

import (
	"context"

	"google.golang.org/grpc"

	"github.com/olezhek28/microservices_course/week_8/circuit_breaker/internal/metric"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)

	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
	}

	return res, err
}
