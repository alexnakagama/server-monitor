package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type ServerRepository interface {
	Create(ctx context.Context, server model.Server) error
}

type ServerService struct {
	repository ServerRepository
}

func (s *ServerService) Create(ctx context.Context, name string, hostname string, os string) error {
	server := model.Server{
		Name:     name,
		Hostname: hostname,
		OS:       os,
	}

	return s.repository.Create(ctx, server)
}
