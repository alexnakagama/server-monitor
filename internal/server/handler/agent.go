package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexnakagama/server-monitor/internal/server/service"
)

type AgentHandler struct {
	service *service.AgentService
}

func NewAgentHandler(service *service.AgentService) *AgentHandler {
	return &AgentHandler{
		service: service,
	}
}

type CreateAgentRequest struct {
	Name string `json:"name"`
}

func (h *AgentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	serverID, err := strconv.Atoi(r.PathValue("serverID"))
	if err != nil {
		http.Error(w, "invalid server id", http.StatusBadRequest)
		return
	}

	var req CreateAgentRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.service.Create(r.Context(), serverID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (h *AgentHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {}
