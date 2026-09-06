package handler

import "github.com/alexnakagama/server-monitor/internal/server/service"

type AgentHandler struct {
	service *service.AgentService
}

func NewAgentHandler(s *service.AgentService) *AgentHandler {
	return &AgentHandler{
		service: s,
	}
}
