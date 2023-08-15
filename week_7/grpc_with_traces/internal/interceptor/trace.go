package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

func TracesInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return handler(ctx, req)
}
