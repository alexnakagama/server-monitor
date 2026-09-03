package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
	"github.com/alexnakagama/server-monitor/internal/server/service"
)

type ServerHandler struct {
	service *service.ServerService
}

func NewServerHandler(service *service.ServerService) *ServerHandler {
	return &ServerHandler{
		service: service,
	}
}

type CreateServerRequest struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

func (h *ServerHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateServerRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.Create(
		r.Context(),
		req.Name,
		req.Hostname,
		req.OS,
	)
	if err != nil {
		http.Error(w, "failed to create a server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type ServerResponse struct {
	Name      string    `json:"name"`
	Hostname  string    `json:"hostname"`
	OS        string    `json:"os"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *ServerHandler) HandleGetByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	server, err := h.service.GetByName(r.Context(), name)
	if err != nil {
		if errors.Is(err, errors_custom.ErrServerNotFound) {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}

		http.Error(w, "server not found", http.StatusNotFound)
		return
	}

	response := ServerResponse{
		Name:      server.Name,
		Hostname:  server.Hostname,
		OS:        server.OS,
		CreatedAt: server.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ServerHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	servers, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to get servers", http.StatusInternalServerError)
		return
	}

	responses := make([]ServerResponse, 0, len(servers))

	for _, server := range servers {
		responses = append(responses, ServerResponse{
			Name:      server.Name,
			Hostname:  server.Hostname,
			OS:        server.OS,
			CreatedAt: server.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *ServerHandler) HandleGetByOS(w http.ResponseWriter, r *http.Request) {
	os := r.PathValue("os")

	servers, err := h.service.GetByOS(r.Context(), os)
	if err != nil {
		http.Error(w, "failed to get servers", http.StatusInternalServerError)
		return
	}

	responses := make([]ServerResponse, 0, len(servers))

	for _, server := range servers {
		responses = append(responses, ServerResponse{
			Name:      server.Name,
			Hostname:  server.Hostname,
			OS:        server.OS,
			CreatedAt: server.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *ServerHandler) HandleGetByHostname(w http.ResponseWriter, r *http.Request) {
	hostname := r.PathValue("hostname")

	server, err := h.service.GetByHostname(r.Context(), hostname)
	if err != nil {
		http.Error(w, "failed to get server", http.StatusInternalServerError)
		return
	}

	response := ServerResponse{
		Name:      server.Name,
		Hostname:  server.Hostname,
		OS:        server.OS,
		CreatedAt: server.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ServerHandler) HandleDeleteByHostname(w http.ResponseWriter, r *http.Request) {
	hostname := r.PathValue("hostname")

	err := h.service.DeleteByHostname(r.Context(), hostname)
	if err != nil {
		if errors.Is(err, errors_custom.ErrServerNotFound) {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to delete server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateServerRequest struct {
	Name string `json:"name"`
	OS   string `json:"os"`
}

func (h *ServerHandler) HandleUpdateByHostname(w http.ResponseWriter, r *http.Request) {
	hostname := r.PathValue("hostname")

	var req UpdateServerRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateByHostname(
		r.Context(),
		hostname,
		req.Name,
		req.OS,
	)
	if err != nil {
		if errors.Is(err, errors_custom.ErrServerNotFound) {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to update server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
