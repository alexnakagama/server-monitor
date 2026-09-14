package tui

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
)

type DashboardModel struct {
	client         *monitor.Client
	servers        []model.Server
	selectedServer int
	err            error
}

type serversLoadedMessage struct {
	servers []model.Server
	err     error
}

func (m *DashboardModel) loadServers() tea.Cmd {
	return func() tea.Msg {
		servers, err := m.client.GetServers(context.Background())

		return serversLoadedMessage{
			servers: servers,
			err:     err,
		}
	}
}

func NewDashBoardModel(client *monitor.Client) DashboardModel {
	return DashboardModel{
		client: client,
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	return m.loadServers()
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

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.selectedServer > 0 {
				m.selectedServer--
			}

		case "down":
			if m.selectedServer < len(m.servers)-1 {
				m.selectedServer++
			}

		case "q":
			return m, func() tea.Msg {
				return quit()
			}
		}
	}

	return m, nil
}

func (m *DashboardModel) View() string {
	view := "Server Monitor\n\n"
	view += "Dashboard\n\n"

	if m.err != nil {
		view += "Error: " + m.err.Error() + "\n"
		return view
	}

	if len(m.servers) == 0 {
		view += "No servers found.\n"
		return view
	}

	view += "Servers:\n\n"

	for i, server := range m.servers {
		if i == m.selectedServer {
			view += "> " + server.Name + "\n"
		} else {
			view += "  " + server.Name + "\n"
		}
	}

	view += "\nUse ↑/↓ to select a server\n"
	view += "Press Enter to open\n"
	view += "Press q to quit\n"

	return view
}
