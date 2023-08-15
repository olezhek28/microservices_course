package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	descAccess "github.com/olezhek28/microservices_course/week_6/jwt/pkg/access_v1"
	descAuth "github.com/olezhek28/microservices_course/week_6/jwt/pkg/auth_v1"
)

const grpcPort = 50051

type serverAuth struct {
	descAuth.UnimplementedAuthV1Server
}

func (s *serverAuth) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRefreshToken not implemented")
}

func (s *serverAuth) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccessToken not implemented")
}

type serverAccess struct {
	descAccess.UnimplementedAccessV1Server
}

func (s *serverAccess) Check(ctx context.Context, req *descAccess.CheckRequest) (*descAccess.CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	descAuth.RegisterAuthV1Server(s, &serverAuth{})
	descAccess.RegisterAccessV1Server(s, &serverAccess{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
