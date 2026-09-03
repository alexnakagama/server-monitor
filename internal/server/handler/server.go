package handler

import "github.com/alexnakagama/server-monitor/internal/server/service"

type ServerHandler struct {
	service *service.ServerService
}

func NewServerHandler(service *service.ServerService) *ServerHandler {
	return &ServerHandler{
		service: service,
	}
}
