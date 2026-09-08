package handler

import (
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/server/service"
)

type MetricHandler struct {
	service *service.MetricService
}

func NewMetricHandler(service *service.MetricService) *MetricHandler {
	return &MetricHandler{
		service: service,
	}
}

func (h *MetricHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {}
