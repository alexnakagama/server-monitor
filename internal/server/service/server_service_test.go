package service

import (
	"context"
	"testing"

	"github.com/alexnakagama/server-monitor/internal/model"
)

var _ ServerRepository = (*ServerRepositoryMock)(nil)

type ServerRepositoryMock struct {
	create               func(ctx context.Context, server model.Server) error
	getByName            func(ctx context.Context, name string) (model.Server, error)
	getAll               func(ctx context.Context) ([]model.Server, error)
	getByOS              func(ctx context.Context, os string) ([]model.Server, error)
	getByHostname        func(ctx context.Context, hostname string) (model.Server, error)
	deleteByHostname     func(ctx context.Context, hostname string) error
	updateByHostname     func(ctx context.Context, hostname string, name string, os string) error
	updateNameByHostname func(ctx context.Context, name string, hostname string) error
	updateOSByHostname   func(ctx context.Context, os string, hostname string) error
}

func (m *ServerRepositoryMock) Create(ctx context.Context, server model.Server) error {
	return m.create(ctx, server)
}

func (m *ServerRepositoryMock) GetByName(ctx context.Context, name string) (model.Server, error) {
	return m.getByName(ctx, name)
}

func (m *ServerRepositoryMock) GetAll(ctx context.Context) ([]model.Server, error) {
	return m.getAll(ctx)
}

func (m *ServerRepositoryMock) GetByOS(ctx context.Context, os string) ([]model.Server, error) {
	return m.getByOS(ctx, os)
}

func (m *ServerRepositoryMock) GetByHostname(ctx context.Context, hostname string) (model.Server, error) {
	return m.getByHostname(ctx, hostname)
}

func (m *ServerRepositoryMock) DeleteByHostname(ctx context.Context, hostname string) error {
	return m.deleteByHostname(ctx, hostname)
}

func (m *ServerRepositoryMock) UpdateByHostname(ctx context.Context, hostname string, name string, os string) error {
	return m.updateByHostname(ctx, hostname, name, os)
}

func (m *ServerRepositoryMock) UpdateNameByHostname(ctx context.Context, name string, hostname string) error {
	return m.updateNameByHostname(ctx, name, hostname)
}

func (m *ServerRepositoryMock) UpdateOSByHostname(ctx context.Context, os string, hostname string) error {
	return m.updateOSByHostname(ctx, os, hostname)
}

func TestServerService_Create(t *testing.T) {
	var createdServer model.Server

	repository := &ServerRepositoryMock{
		create: func(ctx context.Context, server model.Server) error {
			createdServer = server
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.Create(
		context.Background(),
		"web-server",
		"web-01.example.com",
		"ubuntu",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if createdServer.Name != "web-server" {
		t.Errorf("expected name web-server, got: %s", createdServer.Name)
	}

	if createdServer.Hostname != "web-01.example.com" {
		t.Errorf("expected hostname web-01.example.com, got: %s", createdServer.Hostname)
	}

	if createdServer.OS != "ubuntu" {
		t.Errorf("expected os ubuntu, got: %s", createdServer.OS)
	}
}

func TestServerService_Create_ValidationError(t *testing.T) {
	createCalled := false

	repository := &ServerRepositoryMock{
		create: func(ctx context.Context, server model.Server) error {
			createCalled = true
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.Create(
		context.Background(),
		"",
		"",
		"",
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if createCalled {
		t.Errorf("expected Create to not be called")
	}
}

func TestServerService_GetByName(t *testing.T) {
	server := model.Server{
		ID:       1,
		Name:     "web-server",
		Hostname: "web-01.example.com",
		OS:       "ubuntu",
	}

	repository := &ServerRepositoryMock{
		getByName: func(ctx context.Context, name string) (model.Server, error) {
			return server, nil
		},
	}

	service := NewServerService(repository)

	serverFound, err := service.GetByName(context.Background(), "web-server")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if serverFound.ID != server.ID {
		t.Errorf("expected id: %d, got: %d", server.ID, serverFound.ID)
	}

	if serverFound.Name != server.Name {
		t.Errorf("expected name: %s, got: %s", server.Name, serverFound.Name)
	}

	if serverFound.Hostname != server.Hostname {
		t.Errorf("expected hostname: %s, got: %s", server.Hostname, serverFound.Hostname)
	}

	if serverFound.OS != server.OS {
		t.Errorf("expected os: %s, got: %s", server.OS, serverFound.OS)
	}
}

func TestServerService_GetAll(t *testing.T) {
	servers := []model.Server{
		{
			ID:       1,
			Name:     "web-server",
			Hostname: "web-01.example.com",
			OS:       "ubuntu",
		},
		{
			ID:       2,
			Name:     "db-server",
			Hostname: "db-01.example.com",
			OS:       "ubuntu",
		},
	}

	repository := &ServerRepositoryMock{
		getAll: func(ctx context.Context) ([]model.Server, error) {
			return servers, nil
		},
	}

	service := NewServerService(repository)

	serversFound, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(serversFound) != 2 {
		t.Fatalf("expected 2 servers, got: %d", len(serversFound))
	}

	if serversFound[0].ID != servers[0].ID {
		t.Errorf("expected server id: %d, got: %d", servers[0].ID, serversFound[0].ID)
	}

	if serversFound[1].ID != servers[1].ID {
		t.Errorf("expected server id: %d, got: %d", servers[1].ID, serversFound[1].ID)
	}
}

func TestServerService_GetByOS(t *testing.T) {
	servers := []model.Server{
		{
			ID:       1,
			Name:     "web-server",
			Hostname: "web-01.example.com",
			OS:       "ubuntu",
		},
		{
			ID:       2,
			Name:     "web-server-2",
			Hostname: "web-02.example.com",
			OS:       "ubuntu",
		},
	}

	var receivedOS string

	repository := &ServerRepositoryMock{
		getByOS: func(ctx context.Context, os string) ([]model.Server, error) {
			receivedOS = os
			return servers, nil
		},
	}

	service := NewServerService(repository)

	serversFound, err := service.GetByOS(context.Background(), "ubuntu")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if receivedOS != "ubuntu" {
		t.Errorf("expected os ubuntu, got: %s", receivedOS)
	}

	if len(serversFound) != 2 {
		t.Fatalf("expected 2 servers, got: %d", len(serversFound))
	}

	if serversFound[0].ID != servers[0].ID {
		t.Errorf("expected server id: %d, got: %d", servers[0].ID, serversFound[0].ID)
	}

	if serversFound[1].ID != servers[1].ID {
		t.Errorf("expected server id: %d, got: %d", servers[1].ID, serversFound[1].ID)
	}
}

func TestServerService_GetByHostname(t *testing.T) {
	server := model.Server{
		ID:       1,
		Name:     "web-server",
		Hostname: "web-01.example.com",
		OS:       "ubuntu",
	}

	repository := &ServerRepositoryMock{
		getByHostname: func(ctx context.Context, hostname string) (model.Server, error) {
			return server, nil
		},
	}

	service := NewServerService(repository)

	serverFound, err := service.GetByHostname(context.Background(), "web-01.example.com")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if serverFound.ID != server.ID {
		t.Errorf("expected id: %d, got: %d", server.ID, serverFound.ID)
	}

	if serverFound.Name != server.Name {
		t.Errorf("expected name: %s, got: %s", server.Name, serverFound.Name)
	}

	if serverFound.Hostname != server.Hostname {
		t.Errorf("expected hostname: %s, got: %s", server.Hostname, serverFound.Hostname)
	}

	if serverFound.OS != server.OS {
		t.Errorf("expected os: %s, got: %s", server.OS, serverFound.OS)
	}
}

func TestServerService_GetByHostname_ValidationError(t *testing.T) {
	getByHostnameCalled := false

	repository := &ServerRepositoryMock{
		getByHostname: func(ctx context.Context, hostname string) (model.Server, error) {
			getByHostnameCalled = true
			return model.Server{}, nil
		},
	}

	service := NewServerService(repository)

	_, err := service.GetByHostname(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if getByHostnameCalled {
		t.Errorf("expected GetByHostname to not be called")
	}
}

func TestServerService_DeleteByHostname(t *testing.T) {
	var deletedHostname string

	repository := &ServerRepositoryMock{
		deleteByHostname: func(ctx context.Context, hostname string) error {
			deletedHostname = hostname
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.DeleteByHostname(context.Background(), "web-01.example.com")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if deletedHostname != "web-01.example.com" {
		t.Errorf("expected hostname web-01.example.com, got: %s", deletedHostname)
	}
}

func TestServerService_DeleteByHostname_ValidationError(t *testing.T) {
	deleteByHostnameCalled := false

	repository := &ServerRepositoryMock{
		deleteByHostname: func(ctx context.Context, hostname string) error {
			deleteByHostnameCalled = true
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.DeleteByHostname(context.Background(), "   ")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if deleteByHostnameCalled {
		t.Errorf("expected DeleteByHostname to not be called")
	}
}

func TestServerService_UpdateByHostname(t *testing.T) {
	var receivedHostname string
	var receivedName string
	var receivedOS string

	repository := &ServerRepositoryMock{
		updateByHostname: func(ctx context.Context, hostname string, name string, os string) error {
			receivedHostname = hostname
			receivedName = name
			receivedOS = os
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.UpdateByHostname(
		context.Background(),
		"web-01.example.com",
		"web-server",
		"ubuntu",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if receivedHostname != "web-01.example.com" {
		t.Errorf("expected hostname web-01.example.com, got: %s", receivedHostname)
	}

	if receivedName != "web-server" {
		t.Errorf("expected name web-server, got: %s", receivedName)
	}

	if receivedOS != "ubuntu" {
		t.Errorf("expected os ubuntu, got: %s", receivedOS)
	}
}

func TestServerService_UpdateByHostname_ValidationError(t *testing.T) {
	updateByHostnameCalled := false

	repository := &ServerRepositoryMock{
		updateByHostname: func(ctx context.Context, hostname string, name string, os string) error {
			updateByHostnameCalled = true
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.UpdateByHostname(
		context.Background(),
		"",
		"",
		"",
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if updateByHostnameCalled {
		t.Errorf("expected UpdateByHostname to not be called")
	}
}

func TestServerService_UpdateNameByHostname(t *testing.T) {
	var receivedName string
	var receivedHostname string

	repository := &ServerRepositoryMock{
		updateNameByHostname: func(ctx context.Context, name string, hostname string) error {
			receivedName = name
			receivedHostname = hostname
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.UpdateNameByHostname(
		context.Background(),
		"web-server",
		"web-01.example.com",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if receivedName != "web-server" {
		t.Errorf("expected name web-server, got: %s", receivedName)
	}

	if receivedHostname != "web-01.example.com" {
		t.Errorf("expected hostname web-01.example.com, got: %s", receivedHostname)
	}
}

func TestServerService_UpdateNameByHostname_ValidationError(t *testing.T) {
	updateNameByHostnameCalled := false

	repository := &ServerRepositoryMock{
		updateNameByHostname: func(ctx context.Context, name string, hostname string) error {
			updateNameByHostnameCalled = true
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.UpdateNameByHostname(
		context.Background(),
		"",
		"",
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if updateNameByHostnameCalled {
		t.Errorf("expected UpdateNameByHostname to not be called")
	}
}

func TestServerService_UpdateOSByHostname(t *testing.T) {
	var receivedOS string
	var receivedHostname string

	repository := &ServerRepositoryMock{
		updateOSByHostname: func(ctx context.Context, os string, hostname string) error {
			receivedOS = os
			receivedHostname = hostname
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.UpdateOSByHostname(
		context.Background(),
		"ubuntu",
		"web-01.example.com",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if receivedOS != "ubuntu" {
		t.Errorf("expected os ubuntu, got: %s", receivedOS)
	}

	if receivedHostname != "web-01.example.com" {
		t.Errorf("expected hostname web-01.example.com, got: %s", receivedHostname)
	}
}

func TestServerService_UpdateOSByHostname_ValidationError(t *testing.T) {
	updateOSByHostnameCalled := false

	repository := &ServerRepositoryMock{
		updateOSByHostname: func(ctx context.Context, os string, hostname string) error {
			updateOSByHostnameCalled = true
			return nil
		},
	}

	service := NewServerService(repository)

	err := service.UpdateOSByHostname(
		context.Background(),
		"",
		"",
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if updateOSByHostnameCalled {
		t.Errorf("expected UpdateOSByHostname to not be called")
	}
}
