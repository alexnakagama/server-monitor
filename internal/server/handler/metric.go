package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alexnakagama/server-monitor/internal/auth"
	"github.com/alexnakagama/server-monitor/internal/model"
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

type MetricResponse struct {
	ID             int       `json:"id"`
	ServerID       int       `json:"server_id"`
	CPUUsage       float64   `json:"cpu_usage"`
	MemoryUsage    float64   `json:"memory_usage"`
	DiskUsage      float64   `json:"disk_usage"`
	NetworkReceive uint64    `json:"network_receive"`
	NetworkSent    uint64    `json:"network_sent"`
	Timestamp      time.Time `json:"timestamp"`
}

func (h *MetricHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	metricID, err := strconv.Atoi(r.PathValue("metricID"))
	if err != nil {
		http.Error(w, "invalid metric id", http.StatusBadRequest)
		return
	}

	metric, err := h.service.GetByID(r.Context(), metricID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := MetricResponse{
		ID:             metric.ID,
		ServerID:       metric.ServerID,
		CPUUsage:       metric.CPUUsage,
		MemoryUsage:    metric.MemoryUsage,
		DiskUsage:      metric.DiskUsage,
		NetworkReceive: metric.NetworkReceive,
		NetworkSent:    metric.NetworkSent,
		Timestamp:      metric.Timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (h *MetricHandler) HandleGetByServerID(w http.ResponseWriter, r *http.Request) {
	serverIDStr := r.PathValue("serverID")

	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		http.Error(w, "invalid server id", http.StatusBadRequest)
		return
	}

	limit := 100

	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		if limit < 1 || limit > 1000 {
			http.Error(w, "limit must be between 1 and 1000", http.StatusBadRequest)
			return
		}
	}

	var from *time.Time

	fromParam := r.URL.Query().Get("from")
	if fromParam != "" {
		parsedFrom, err := time.Parse(time.RFC3339, fromParam)
		if err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
			return
		}

		from = &parsedFrom
	}

	var to *time.Time

	toParam := r.URL.Query().Get("to")
	if toParam != "" {
		parsedTo, err := time.Parse(time.RFC3339, toParam)
		if err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
			return
		}

		to = &parsedTo
	}

	filters := model.MetricFilters{
		From:  from,
		To:    to,
		Limit: limit,
	}

	metrics, err := h.service.GetByServerID(r.Context(), serverID, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]MetricResponse, 0, len(metrics))

	for _, metric := range metrics {
		response = append(response, MetricResponse{
			ID:             metric.ID,
			ServerID:       metric.ServerID,
			CPUUsage:       metric.CPUUsage,
			MemoryUsage:    metric.MemoryUsage,
			DiskUsage:      metric.DiskUsage,
			NetworkReceive: metric.NetworkReceive,
			NetworkSent:    metric.NetworkSent,
			Timestamp:      metric.Timestamp,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
