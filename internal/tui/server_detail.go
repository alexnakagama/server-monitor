package tui

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

type tickMessage struct{}

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

func (m *ServerDetailModel) tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMessage{}
	})
}

func NewServerDetailModel(client *monitor.Client, server model.Server) ServerDetailModel {
	return ServerDetailModel{
		client: client,
		server: server,
	}
}

func (m *ServerDetailModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadMetrics(),
		m.tick(),
	)
}

func (m *ServerDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, quit()

		case "esc":
			dashboard := NewDashBoardModel(m.client)
			return &dashboard, dashboard.Init()
		}

	case metricsLoadedMessage:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.metrics = msg.metrics

	case tickMessage:
		return m, tea.Batch(
			m.loadMetrics(),
			m.tick(),
		)
	}

	return m, nil
}

func (m *ServerDetailModel) View() string {
	view := "Server Monitor\n\n"
	view += "Server Details\n\n"

	view += "ID: " + strconv.Itoa(m.server.ID) + "\n"
	view += "Name: " + m.server.Name + "\n"
	view += "Hostname: " + m.server.Hostname + "\n"
	view += "OS: " + m.server.OS + "\n"

	if m.err != nil {
		view += "\nError: " + m.err.Error() + "\n"
		return view
	}

	if len(m.metrics) == 0 {
		view += "\nNo metrics found.\n"
		return view
	}

	metric := m.metrics[0]

	view += "\nMetrics\n\n"
	view += fmt.Sprintf("CPU: %.2f%%\n", metric.CPUUsage)
	view += fmt.Sprintf("Memory: %.2f%%\n", metric.MemoryUsage)
	view += fmt.Sprintf("Disk: %.2f%%\n", metric.DiskUsage)
	view += fmt.Sprintf("Network RX: %d\n", metric.NetworkReceive)
	view += fmt.Sprintf("Network TX: %d\n", metric.NetworkSent)
	view += fmt.Sprintf("Last update: %s\n", metric.Timestamp.Format("15:04:05"))

	view += "\nPress ESC to go back\n"
	view += "\nPress q to quit\n"

	return view
}
