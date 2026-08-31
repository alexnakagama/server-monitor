package handler

import (
	"context"
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/server/service"
)

type UserRespository interface {
	Create(ctx context.Context, user model.User) error
	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user model.User) error
}

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {}
