package handler

import (
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/server/service"
)

type MetricHandler struct {
	service *service.ServerService
}

func NewMetricHandler(service *service.ServerService) *MetricHandler {
	return &MetricHandler{
		service: service,
	}
}

func (h *MetricHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {}
