package tui

import (
	"strings"

	"github.com/alexnakagama/server-monitor/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	switch screen := m.screen.(type) {
	case *LoginModel:
		loginView := screen.View()

		header := lipgloss.NewStyle().
			Width(m.width - 2).
			Render(
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					titleStyle.Render("SERVER MONITOR"),
					lipgloss.PlaceHorizontal(
						m.width-2-lipgloss.Width("SERVER MONITOR"),
						lipgloss.Right,
						helpStyle.Render("v1.0"),
					),
				),
			)

		content := lipgloss.Place(
			m.width,
			m.height-5,
			lipgloss.Center,
			lipgloss.Center,
			loginView,
		)

		divider := strings.Repeat("─", m.width-2)

		footer := lipgloss.PlaceHorizontal(
			m.width-2,
			lipgloss.Right,
			helpStyle.Render("[Q] Quit"),
		)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			divider,
			content,
			divider,
			footer,
		)
	}

	view := m.screen.View()

	container := containerStyle.
		Width(m.width - 2).
		Height(m.height - 2)

	return container.Render(view)
}
