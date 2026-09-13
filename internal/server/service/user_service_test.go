package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type UserRepositoryMock struct {
	getByUsername func(ctx context.Context, username string) (model.User, error)
	getByEmail    func(ctx context.Context, email string) (model.User, error)
	create        func(ctx context.Context, user model.User) error
}
