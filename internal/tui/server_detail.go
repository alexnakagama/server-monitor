package tui

import (
	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/monitor"
)

type ServerDetailModel struct {
	client *monitor.Client
	server model.Server
}
