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

type CreateMetricRequest struct {
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsage    float64 `json:"memory_usage"`
	DiskUsage      float64 `json:"disk_usage"`
	NetworkReceive uint64  `json:"network_receive"`
	NetworkSent    uint64  `json:"network_sent"`
}

func (h *MetricHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
}
