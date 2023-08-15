package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/natefinch/lumberjack"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/olezhek28/microservices_course/week_7/grpc_with_traces/internal/client/rpc"
	otherService "github.com/olezhek28/microservices_course/week_7/grpc_with_traces/internal/client/rpc/other_service"
	"github.com/olezhek28/microservices_course/week_7/grpc_with_traces/internal/interceptor"
	"github.com/olezhek28/microservices_course/week_7/grpc_with_traces/internal/logger"
	"github.com/olezhek28/microservices_course/week_7/grpc_with_traces/internal/tracing"
	desc "github.com/olezhek28/microservices_course/week_7/grpc_with_traces/pkg/note_v1"
	descOther "github.com/olezhek28/microservices_course/week_7/grpc_with_traces/pkg/other_note_v1"
)

var logLevel = flag.String("l", "info", "log level")

const (
	grpcPort         = 50051
	otherServicePort = 50052
	serviceName      = "test-service"
)

type server struct {
	desc.UnimplementedNoteV1Server

	otherServiceClient rpc.OtherServiceClient
}

// Get ...
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	if req.GetId() == 0 {
		return nil, errors.Errorf("id is empty")
	}

	// rand.Intn(max - min) + min
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	span, ctx := opentracing.StartSpanFromContext(ctx, "get note")
	defer span.Finish()

	span.SetTag("id", req.GetId())

	note, err := s.otherServiceClient.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.WithMessage(err, "getting note")
	}

	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: note.ID,
			Info: &desc.NoteInfo{
				Title:   note.Info.Title,
				Content: note.Info.Content,
			},
			CreatedAt: timestamppb.New(note.CreatedAt),
			UpdatedAt: updatedAt,
		},
	}, nil
}

func main() {
	flag.Parse()

	logger.Init(getCore(getAtomicLevel()))
	tracing.Init(logger.Logger(), serviceName)

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", otherServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatalf("failed to dial GRPC client: %v", err)
	}

	otherServiceClient := otherService.New(descOther.NewOtherNoteV1Client(conn))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.ServerTracingInterceptor,
		),
	)
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{otherServiceClient: otherServiceClient})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
