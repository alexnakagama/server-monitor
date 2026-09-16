package tui

import (
	"context"
	"fmt"
	"strings"
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
	width          int
	height         int
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

func (m *DashboardModel) filteredServers() ([]model.Server, []string) {
	searchTerm := strings.ToLower(strings.TrimSpace(m.search.Value()))

	if searchTerm == "" {
		return m.servers, m.statuses
	}

	var servers []model.Server
	var statuses []string

	for i, server := range m.servers {
		if !strings.Contains(
			strings.ToLower(server.Name),
			searchTerm,
		) {
			continue
		}

		servers = append(servers, server)
		statuses = append(statuses, m.statuses[i])
	}

	return servers, statuses
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

func (m *DashboardModel) SetSize(width, height int) {
	m.width = width
	m.height = height
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
	case tea.WindowSizeMsg:
		fmt.Printf("DASHBOARD WIDTH: %d\n", msg.Width)

		m.width = msg.Width
		m.height = msg.Height

		return m, nil

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

			case "up":
				if m.selectedServer > 0 {
					m.selectedServer--
				}

				return m, nil

			case "down":
				servers, _ := m.filteredServers()

				if m.selectedServer < len(servers)-1 {
					m.selectedServer++
				}

				return m, nil

			case "enter":
				servers, _ := m.filteredServers()

				if len(servers) == 0 {
					return m, nil
				}

				server := servers[m.selectedServer]

				detail := NewServerDetailModel(m.client, server)
				detail.SetSize(m.width, m.height)

				return &detail, detail.Init()
			}

			var cmd tea.Cmd

			m.search, cmd = m.search.Update(msg)
			m.selectedServer = 0

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

	servers, statuses := m.filteredServers()

	if len(servers) == 0 {
		view += "No servers match your search.\n"
	} else {
		view += sectionStyle.Render("SERVERS") + "\n\n"

		view += prefixStyle.Render("  ") +
			nameColumnStyle.Bold(true).Render("NAME") +
			osColumnStyle.Bold(true).Render("OS") +
			headerStyle.Render("STATUS") +
			"\n"

		for i, server := range servers {
			prefix := "  "

			if i == m.selectedServer {
				prefix = "> "
			}

			name := nameColumnStyle.Render(server.Name)
			os := osColumnStyle.Render(server.OS)

			status := statuses[i]

			statusView := offlineStyle.Render(status)

			if status == "ONLINE" {
				statusView = onlineStyle.Render(status)
			}

			row := prefixStyle.Render(prefix) + name + os + statusView

			view += row + "\n"
		}
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

	container := containerStyle.Width(m.width - 4)

	return container.Render(view)
}
