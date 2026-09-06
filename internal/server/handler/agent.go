package handler

import (
	"net/http"

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
}

func (h *AgentHandler) HandleCreate(w *http.ResponseWriter, r *http.Request) {}
