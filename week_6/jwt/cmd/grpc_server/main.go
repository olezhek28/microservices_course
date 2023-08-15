package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/olezhek28/microservices_course/week_6/jwt/internal/utils"
	descAccess "github.com/olezhek28/microservices_course/week_6/jwt/pkg/access_v1"
	descAuth "github.com/olezhek28/microservices_course/week_6/jwt/pkg/auth_v1"
)

const (
	grpcPort = 50051

	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)

type serverAuth struct {
	descAuth.UnimplementedAuthV1Server
}

func (s *serverAuth) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	// Лезем в базу или кэш за данными пользователя
	// Сверяем хэши пароля

	refreshToken, err := utils.GenerateToken(req.GetUsername(), []byte(refreshTokenSecretKey), refreshTokenExpiration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &descAuth.LoginResponse{RefreshToken: refreshToken}, nil
}

func (s *serverAuth) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	// Можем слазать в базу или в кэш за доп данными пользователя

	refreshToken, err := utils.GenerateToken(claims.Username, []byte(refreshTokenSecretKey), refreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	return &descAuth.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}

func (s *serverAuth) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	// Можем слазать в базу или в кэш за доп данными пользователя

	accessToken, err := utils.GenerateToken(claims.Username, []byte(accessTokenSecretKey), accessTokenExpiration)
	if err != nil {
		return nil, err
	}

	return &descAuth.GetAccessTokenResponse{AccessToken: accessToken}, nil
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
