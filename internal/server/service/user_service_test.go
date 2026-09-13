package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

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
