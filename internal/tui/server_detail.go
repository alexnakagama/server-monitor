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

func progressBar(value float64, width int) string {}

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
	view := titleStyle.Render("SERVER MONITOR") + "\n\n"

	view += sectionStyle.Render("SERVER") + "\n"

	view += labelStyle.Render("ID") + valueStyle.Render(strconv.Itoa(m.server.ID)) + "\n"
	view += labelStyle.Render("Name") + valueStyle.Render(m.server.Name) + "\n"
	view += labelStyle.Render("Hostname") + valueStyle.Render(m.server.Hostname) + "\n"
	view += labelStyle.Render("OS") + valueStyle.Render(m.server.OS) + "\n"

	if m.err != nil {
		view += "\nError: " + m.err.Error() + "\n"
		return view
	}

	if len(m.metrics) == 0 {
		view += "\nNo metrics found.\n"
		return view
	}

	metric := m.metrics[0]

	view += "\n"
	view += sectionStyle.Render("METRICS") + "\n"

	view += labelStyle.Render("CPU") + fmt.Sprintf("%.2f%%\n", metric.CPUUsage)
	view += labelStyle.Render("Memory") + fmt.Sprintf("%.2f%%\n", metric.MemoryUsage)
	view += labelStyle.Render("Disk") + fmt.Sprintf("%.2f%%\n", metric.DiskUsage)

	view += "\n"
	view += sectionStyle.Render("NETWORK") + "\n"

	view += labelStyle.Render("RX") + fmt.Sprintf("%d\n", metric.NetworkReceive)
	view += labelStyle.Render("TX") + fmt.Sprintf("%d\n", metric.NetworkSent)

	view += "\n"
	view += helpStyle.Render(
		"Last update: "+metric.Timestamp.Format("15:04:05"),
	) + "\n\n"

	view += helpStyle.Render("[Esc] Back   [q] Quit") + "\n"

	return containerStyle.Render(view)
}
