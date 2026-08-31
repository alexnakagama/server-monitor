package handler

import (
	"context"
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type UserRespository interface {
	Create(ctx context.Context, user model.User) error
	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user model.User) error
}

type UserHandler struct {
	repository UserRespository
}

func NewUserHandler(repo UserRespository) *UserHandler {
	return &UserHandler{
		repository: repo,
	}
}

func (h *UserHandler) HandlerCreate(w http.ResponseWriter, r *http.Request) {}
