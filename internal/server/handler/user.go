package handler

import (
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/server/service"
)

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
