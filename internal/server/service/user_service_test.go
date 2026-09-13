package service

import (
	"context"
	"testing"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

var _ UserRepository = (*UserRepositoryMock)(nil)

type UserRepositoryMock struct {
	getByUsername  func(ctx context.Context, username string) (model.User, error)
	getByEmail     func(ctx context.Context, email string) (model.User, error)
	create         func(ctx context.Context, user model.User) error
	delete         func(ctx context.Context, id int) error
	update         func(ctx context.Context, user model.User) error
	getByID        func(ctx context.Context, id int) (model.User, error)
	updatePassword func(ctx context.Context, id int, hash string) error
}

func (m *UserRepositoryMock) GetByUsername(ctx context.Context, username string) (model.User, error) {
	return m.getByUsername(ctx, username)
}

func (m *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return m.getByEmail(ctx, email)
}

func (m *UserRepositoryMock) Create(ctx context.Context, user model.User) error {
	return m.create(ctx, user)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, id int) error {
	return m.delete(ctx, id)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user model.User) error {
	return m.update(ctx, user)
}

func (m *UserRepositoryMock) GetByID(ctx context.Context, id int) (model.User, error) {
	return m.getByID(ctx, id)
}

func (m *UserRepositoryMock) UpdatePassword(ctx context.Context, id int, hash string) error {
	return m.updatePassword(ctx, id, hash)
}

func TestUserService_Register(t *testing.T) {
	repository := &UserRepositoryMock{
		getByUsername: func(ctx context.Context, username string) (model.User, error) {
			return model.User{}, errors_custom.ErrUserNotFound
		},
		getByEmail: func(ctx context.Context, email string) (model.User, error) {
			return model.User{}, errors_custom.ErrUserNotFound
		},
		create: func(ctx context.Context, user model.User) error {
			return nil
		},
	}

	service := NewUserService(repository, nil)

	err := service.Register(
		context.Background(),
		"alex123",
		"alex@example.gmail.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func (m *UserRepositoryMock) TestUserService_Login(t *testing.T) {}
