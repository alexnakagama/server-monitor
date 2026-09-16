package tui

import (
	"github.com/alexnakagama/server-monitor/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	width  int
	height int
	screen tea.Model
}

type dashboardMessage struct {
	client *monitor.Client
}

func NewAppModel(screen tea.Model) AppModel {
	return AppModel{
		screen: screen,
	}
}

func (m *AppModel) Init() tea.Cmd {
	return m.screen.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil

	case dashboardMessage:
		dashboard := NewDashBoardModel(msg.client)

		m.screen = &dashboard

		return m, dashboard.Init()
	}

	var cmd tea.Cmd

	m.screen, cmd = m.screen.Update(msg)

	return m, cmd
}

func (m *AppModel) View() string {
	view := m.screen.View()

	container := containerStyle.
		Width(m.width - 4).
		Height(m.height - 2)

	return container.Render(view)
}
