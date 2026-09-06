package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/auth"
	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
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

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.Register(
		r.Context(),
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *UserHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (h *UserHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateProfile(
		r.Context(),
		userID,
		req.Username,
		req.Email,
	)

	if err != nil {
		switch {
		case errors.Is(err, errors_custom.ErrUsernameRequired),
			errors.Is(err, errors_custom.ErrUsernameTooShort),
			errors.Is(err, errors_custom.ErrUsernameTooLong),
			errors.Is(err, errors_custom.ErrEmailRequired),
			errors.Is(err, errors_custom.ErrEmailTooLong),
			errors.Is(err, errors_custom.ErrInvalidEmail):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, errors_custom.ErrUsernameAlreadyExists),
			errors.Is(err, errors_custom.ErrEmailAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)

		case errors.Is(err, errors_custom.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) HandleDeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.DeleteProfile(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, errors_custom.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *UserHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.ChangePassword(r.Context())
}
