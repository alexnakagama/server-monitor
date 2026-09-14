package tui

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
)

type DashboardModel struct {
	client  *monitor.Client
	servers []model.Server
	err     error
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
		case "q":
			return m, func() tea.Msg {
				return quit()
			}
		}
	}

	return m, nil
}

func (m *DashboardModel) View() string {
	return "Server Monitor\n\n" +
		"Dashboard\n\n" +
		"Logged in successfully.\n\n" +
		"Press q to quit\n"
}
