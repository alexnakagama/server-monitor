package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/monitor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type DashboardModel struct {
	client         *monitor.Client
	servers        []model.Server
	selectedServer int
	err            error
	statuses       []string
	search         textinput.Model
	searching      bool
}

type serversLoadedMessage struct {
	servers  []model.Server
	err      error
	statuses []string
}

func (m *DashboardModel) tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMessage{}
	})
}

func (m *DashboardModel) loadServers() tea.Cmd {
	return func() tea.Msg {
		servers, err := m.client.GetServers(context.Background())
		if err != nil {
			return serversLoadedMessage{
				err: err,
			}
		}

		statuses := make([]string, len(servers))

		for i, server := range servers {
			metric, err := m.client.GetLatestServerMetric(
				context.Background(),
				server.ID,
			)

			if err != nil {
				statuses[i] = "OFFLINE"
				continue
			}

			if time.Since(metric.Timestamp) < 30*time.Second {
				statuses[i] = "ONLINE"
			} else {
				statuses[i] = "OFFLINE"
			}
		}

		return serversLoadedMessage{
			servers:  servers,
			statuses: statuses,
		}
	}
}

func NewDashBoardModel(client *monitor.Client) DashboardModel {
	search := textinput.New()
	search.Placeholder = "Search servers..."
	search.CharLimit = 50

	return DashboardModel{
		client: client,
		search: search,
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadServers(),
		m.tick(),
	)
}

func quit() tea.Cmd {
	return func() tea.Msg {
		return tea.Quit()
	}
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case serversLoadedMessage:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.servers = msg.servers
		m.statuses = msg.statuses

		return m, nil

	case tickMessage:
		return m, tea.Batch(
			m.loadServers(),
			m.tick(),
		)

	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "esc":
				m.search.Reset()
				m.search.Blur()
				m.searching = false
				m.selectedServer = 0

				return m, nil
			}

			var cmd tea.Cmd

			m.search, cmd = m.search.Update(msg)

			return m, cmd
		}

		switch msg.String() {
		case "/":
			m.searching = true
			m.selectedServer = 0

			cmd := m.search.Focus()

			return m, cmd

		case "up":
			if m.selectedServer > 0 {
				m.selectedServer--
			}

		case "down":
			if m.selectedServer < len(m.servers)-1 {
				m.selectedServer++
			}

		case "enter":
			if len(m.servers) == 0 {
				return m, nil
			}

			server := m.servers[m.selectedServer]
			detail := NewServerDetailModel(m.client, server)

			return &detail, detail.Init()

		case "q":
			return m, quit()
		}
	}

	return m, nil
}

func (m *DashboardModel) View() string {
	view := titleStyle.Render("SERVER MONITOR") + "\n\n"

	view += sectionStyle.Render("DASHBOARD") + "\n\n"

	if m.err != nil {
		view += "Error: " + m.err.Error() + "\n"
		return containerStyle.Render(view)
	}

	if m.searching {
		view += m.search.View() + "\n\n"
	}

	if len(m.servers) == 0 {
		view += "No servers found.\n"
		view += "\n"
		view += helpStyle.Render("[/] Search   [q] Quit") + "\n"

		return containerStyle.Render(view)
	}

	view += sectionStyle.Render("SERVERS") + "\n\n"

	view += headerStyle.Render(
		fmt.Sprintf("%-22s %-10s %-10s", "NAME", "OS", "STATUS"),
	) + "\n"

	for i, server := range m.servers {
		prefix := "  "

		if i == m.selectedServer {
			prefix = "> "
		}

		name := nameColumnStyle.Render(server.Name)
		os := osColumnStyle.Render(server.OS)

		status := m.statuses[i]

		statusView := offlineStyle.Render(status)

		if status == "ONLINE" {
			statusView = onlineStyle.Render(status)
		}

		row := prefix + name + os + statusView

		view += row + "\n"
	}

	view += "\n"

	if m.searching {
		view += helpStyle.Render(
			"[Esc] Close search   [↑/↓] Navigate   [Enter] Open",
		) + "\n"
	} else {
		view += helpStyle.Render(
			"[/] Search   [↑/↓] Navigate   [Enter] Open   [q] Quit",
		) + "\n"
	}

	return containerStyle.Render(view)
}
