package tui

import (
	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
)

type ServerDetailModel struct {
	client *monitor.Client
	server model.Server
}

func NewServerDetailModel(client *monitor.Client, server model.Server) ServerDetailModel {
	return ServerDetailModel{
		client: client,
		server: server,
	}
}

func (m *ServerDetailModel) Init() tea.Cmd {
	return nil
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
