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
	GetByHostname(ctx context.Context, hostname string) (model.Server, error)
	DeleteByHostname(ctx context.Context, hostname string) error
	UpdateByHostname(ctx context.Context, hostname string, name string, os string) error
	UpdateNameByHostname(ctx context.Context, name string, hostname string) error
	UpdateOSByHostname(ctx context.Context, os string, hostname string) error
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

	err := server.Validate()
	if err != nil {
		return err
	}

	return s.repository.Create(ctx, server)
}

func (s *ServerService) GetByName(ctx context.Context, name string) (model.Server, error) {
	return s.repository.GetByName(ctx, name)
}

func (s *ServerService) GetAll(ctx context.Context) ([]model.Server, error) {
	return s.repository.GetAll(ctx)
}

func (s *ServerService) GetByOS(ctx context.Context, os string) ([]model.Server, error) {
	return s.repository.GetByOS(ctx, os)
}

func (s *ServerService) GetByHostname(ctx context.Context, hostname string) (model.Server, error) {
	err := model.ValidateHostname(hostname)
	if err != nil {
		return model.Server{}, err
	}

	return s.repository.GetByHostname(ctx, hostname)
}

func (s *ServerService) DeleteByHostname(ctx context.Context, hostname string) error {
	err := model.ValidateHostname(hostname)
	if err != nil {
		return err
	}

	return s.repository.DeleteByHostname(ctx, hostname)
}

func (s *ServerService) UpdateByHostname(ctx context.Context, hostname string, name string, os string) error {
	err := model.ValidateHostname(hostname)
	if err != nil {
		return err
	}

	err = model.ValidateName(name)
	if err != nil {
		return err
	}

	err = model.ValidateOS(os)
	if err != nil {
		return err
	}

	return s.repository.UpdateByHostname(ctx, hostname, name, os)
}

func (s *ServerService) UpdateNameByHostname(ctx context.Context, name string, hostname string) error {
	err := model.ValidateHostname(hostname)
	if err != nil {
		return err
	}

	err = model.ValidateName(name)
	if err != nil {
		return err
	}

	return s.repository.UpdateNameByHostname(ctx, name, hostname)
}

func (s *ServerService) UpdateOSByHostname(ctx context.Context, os string, hostname string) error {
	err := model.ValidateOS(os)
	if err != nil {
		return err
	}

	return s.repository.UpdateOSByHostname(ctx, os, hostname)
}
