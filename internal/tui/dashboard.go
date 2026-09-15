package tui

import (
	"context"
	"fmt"

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

	if len(m.servers) == 0 {
		view += "No servers found.\n"
		return containerStyle.Render(view)
	}

	view += sectionStyle.Render("SERVERS") + "\n\n"

	for i, server := range m.servers {
		row := fmt.Sprintf(
			"%s%-20s %-10s",
			"  ",
			server.Name,
			server.OS,
		)

		if i == m.selectedServer {
			row = fmt.Sprintf(
				"%s%-20s %-10s",
				"> ",
				server.Name,
				server.OS,
			)

			row = selectedStyle.Render(row)
		}

		view += row + "\n"
	}

	view += "\n"
	view += helpStyle.Render("[↑/↓] Navigate   [Enter] Open   [q] Quit") + "\n"

	return containerStyle.Render(view)
}
