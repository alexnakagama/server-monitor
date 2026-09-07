package handler

import "github.com/alexnakagama/server-monitor/internal/server/service"

type MetricHandler struct {
	service *service.ServerService
}

func NewMetricHandler(service *service.ServerService) *MetricHandler {}
