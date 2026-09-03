package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type ServerRepository interface {
	Create(ctx context.Context, server model.Server) error
	GetByName(ctx context.Context, name string) (model.Server, error)
	GetAll(ctx context.Context) ([]model.Server, error)
	GetByOS(ctx context.Context, os string) ([]model.Server, error)
}

type ServerService struct {
	repository ServerRepository
}

func NewServerService(repository ServerRepository) *ServerService {
	return &ServerService{
		repository: repository,
	}
}

func (s *ServerService) Create(ctx context.Context, name string, hostname string, os string) error {
	server := model.Server{
		Name:     name,
		Hostname: hostname,
		OS:       os,
	}

	return s.repository.Create(ctx, server)
}

func (s *ServerService) GetByName(ctx context.Context, name string) (model.Server, error) {
	return s.repository.GetByName(ctx, name)
}

func (s *ServerService) GetAll(ctx context.Context) ([]model.Server, error) {
	return s.repository.GetAll(ctx)
}
