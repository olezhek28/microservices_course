package main

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/olezhek28/microservices_course/week_8/rate_limiter/pkg/note_v1"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: noteID})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}

	log.Printf(color.RedString("Note info:\n"), color.GreenString("%+v", r.GetNote()))
}

func clientOpt() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(grpcRetry.UnaryClientInterceptor(
		grpcRetry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
		grpcRetry.WithMax(5),
		grpcRetry.WithBackoff(grpcRetry.BackoffLinear(time.Second)),
	)))
	opts = append(opts, grpc.WithInsecure())

	grpc.DialContext(context.Background(), "localhost:8080", opts...)
}
