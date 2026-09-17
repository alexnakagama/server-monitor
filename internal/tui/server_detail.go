package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
)

type ServerDetailModel struct {
	client        *monitor.Client
	server        model.Server
	metrics       []model.Metric
	err           error
	showHistory   bool
	historyOffset int
	historyRows   int
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

func progressBar(value float64, width int) string {
	filled := int(value / 100 * float64(width))

	if filled > width {
		filled = width
	}

	if filled < 0 {
		filled = 0
	}

	return fmt.Sprintf(
		"%s╺%s",
		strings.Repeat("━", filled),
		strings.Repeat("━", width-filled),
	)
}

func formatBytes(bytes uint64) string {
	const unit = 1024

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	value := float64(bytes)
	units := []string{"KB", "MB", "GB", "TB"}

	for _, unitName := range units {
		value /= unit

		if value < unit {
			return fmt.Sprintf("%.2f %s", value, unitName)
		}
	}

	return fmt.Sprintf("%.2f PB", value)
}

func (m *ServerDetailModel) historyView() string {
	view := titleStyle.Render("SERVER MONITOR") + "\n\n"

	view += sectionStyle.Render("HISTORY") + "\n\n"

	view += headerStyle.Render(
		fmt.Sprintf("%-10s %-10s %-10s %-10s", "TIME", "CPU", "MEMORY", "DISK"),
	) + "\n"

	for _, metric := range m.metrics {
		view += fmt.Sprintf(
			"%-10s %-10.2f %-10.2f %-10.2f\n",
			metric.Timestamp.Format("15:04:05"),
			metric.CPUUsage,
			metric.MemoryUsage,
			metric.DiskUsage,
		)
	}

	view += "\n"
	view += helpStyle.Render("[Esc] Back") + "\n"

	return view
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
			if m.showHistory {
				m.showHistory = false
				return m, nil
			}

			dashboard := NewDashBoardModel(m.client)

			return &dashboard, dashboard.Init()

		case "h":
			m.showHistory = true
			return m, nil
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
	if m.showHistory {
		return m.historyView()
	}

	view := ""

	// SERVER

	status := "OFFLINE"

	if len(m.metrics) > 0 {
		lastMetric := m.metrics[0]

		if time.Since(lastMetric.Timestamp) < 30*time.Second {
			status = "ONLINE"
		}
	}

	statusView := offlineStyle.Render(status)

	if status == "ONLINE" {
		statusView = onlineStyle.Render(status)
	}

	view += valueStyle.Render(m.server.Name) +
		strings.Repeat(" ", 5) +
		statusView +
		"\n\n"

	view += labelStyle.Render("ID") +
		valueStyle.Render(strconv.Itoa(m.server.ID)) + "\n"

	view += labelStyle.Render("OS") +
		valueStyle.Render(m.server.OS) + "\n"

	view += labelStyle.Render("Hostname") +
		valueStyle.Render(m.server.Hostname) + "\n"

	if len(m.metrics) > 0 {
		metric := m.metrics[0]

		view += labelStyle.Render("Last update") +
			valueStyle.Render(metric.Timestamp.Format("15:04:05")) +
			"\n"
	}

	if m.err != nil {
		view += "\n"
		view += errorStyle.Render(
			"Error: "+m.err.Error(),
		) + "\n"

		return view
	}

	if len(m.metrics) == 0 {
		view += "\n"
		view += helpStyle.Render("No metrics found.") + "\n"

		return view
	}

	metric := m.metrics[0]

	// METRICS
	view += "\n"
	view += sectionStyle.Render("METRICS") + "\n"
	view += strings.Repeat("─", 40) + "\n\n"

	view += labelStyle.Render("CPU") +
		progressBar(metric.CPUUsage, 30) +
		fmt.Sprintf(" %.2f%%\n", metric.CPUUsage)

	view += labelStyle.Render("Memory") +
		progressBar(metric.MemoryUsage, 30) +
		fmt.Sprintf(" %.2f%%\n", metric.MemoryUsage)

	view += labelStyle.Render("Disk") +
		progressBar(metric.DiskUsage, 30) +
		fmt.Sprintf(" %.2f%%\n", metric.DiskUsage)

	// NETWORK
	view += "\n"
	view += sectionStyle.Render("NETWORK") + "\n"
	view += strings.Repeat("─", 40) + "\n\n"

	view += labelStyle.Render("RX") +
		valueStyle.Render(formatBytes(metric.NetworkReceive)+"/s") +
		"\n"

	view += labelStyle.Render("TX") +
		valueStyle.Render(formatBytes(metric.NetworkSent)+"/s") +
		"\n"

	return view
}
