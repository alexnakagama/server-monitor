package tui

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
)

type ServerDetailModel struct {
	client  *monitor.Client
	server  model.Server
	metrics []model.Metric
	err     error
}

type metricsLoadedMessage struct {
	metrics []model.Metric
	err     error
}

func (m *ServerDetailModel) loadMetrics() tea.Cmd {
	return func() tea.Msg {
		metrics, err := m.client.GetServerMetrics(
			context.Background(),
			m.server.ID,
		)

		return metricsLoadedMessage{
			metrics: metrics,
			err:     err,
		}
	}
}

func NewServerDetailModel(client *monitor.Client, server model.Server) ServerDetailModel {
	return ServerDetailModel{
		client: client,
		server: server,
	}
}

func (m *ServerDetailModel) Init() tea.Cmd {
	return m.loadMetrics()
}

func (m *ServerDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, quit()
		}
	}

	return m, nil
}

func (m *ServerDetailModel) View() string {
	view := "Server Monitor\n\n"
	view += "Server Details\n\n"

	view += "Name: " + m.server.Name + "\n"
	view += "Hostname: " + m.server.Hostname + "\n"
	view += "OS: " + m.server.OS + "\n"

	view += "\nPress q to quit\n"

	return view
}
