package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user model.User) error
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}
