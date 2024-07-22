package app

import (
	"context"
	"log"

	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/api/note"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/config"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/config/env"
	"github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/service"
	noteService "github.com/olezhek28/microservices_course/week_5/clean_kafka/internal/service/note"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig

	noteService service.NoteService

	noteImpl *note.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService()
	}

	return s.noteService
}

func (s *serviceProvider) NoteImpl(ctx context.Context) *note.Implementation {
	if s.noteImpl == nil {
		s.noteImpl = note.NewImplementation(s.NoteService(ctx))
	}

	return s.noteImpl
}
