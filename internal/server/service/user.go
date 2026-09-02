package service

import (
	"context"
	"errors"

	"github.com/alexnakagama/server-monitor/internal/model"
	apperrors "github.com/alexnakagama/server-monitor/internal/server/repository"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
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

func (s *UserService) Register(ctx context.Context, username string, email string, password string) error {
	_, err := s.repository.GetByUsername(ctx, username)
	if err != nil {
		return errors.New("username already exists")
	}

	if !errors.Is(err, apperrors.ErrUserNotFound) {
		return err
	}

	_, err = s.repository.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("email already exists")
	}

	if !errors.Is(err, apperrors.ErrUserNotFound) {
		return err
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := model.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}

	return s.repository.Create(ctx, user)
}
