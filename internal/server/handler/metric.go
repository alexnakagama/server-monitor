package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/auth"
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
	serverID, ok := auth.ServerIDFromContext(r.Context())
	if !ok {
		http.Error(w, "server id not found", http.StatusUnauthorized)
		return
	}

	var req CreateMetricRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := service.MetricInput{
		ServerID:       serverID,
		CPUUsage:       req.CPUUsage,
		MemoryUsage:    req.MemoryUsage,
		DiskUsage:      req.DiskUsage,
		NetworkReceive: req.NetworkReceive,
		NetworkSent:    req.NetworkSent,
	}

	err = h.service.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type MetricResponse struct{}

func (h *MetricHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
}
