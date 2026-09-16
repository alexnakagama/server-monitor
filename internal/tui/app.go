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

		headerTitle := titleStyle.Render("SERVER MONITOR")
		headerVersion := helpStyle.Render("v1.0")

		header := lipgloss.JoinHorizontal(
			lipgloss.Top,
			headerTitle,
			lipgloss.NewStyle().
				Width(m.width-2-lipgloss.Width(headerTitle)).
				Align(lipgloss.Right).
				Render(headerVersion),
		)

		divider := strings.Repeat("─", m.width-6)

		content := lipgloss.Place(
			m.width,
			m.height-5,
			lipgloss.Center,
			lipgloss.Center,
			loginView,
		)

		footer := lipgloss.NewStyle().
			Width(m.width - 2).
			Align(lipgloss.Center).
			Render(helpStyle.Render("[Esc] Quit"))

		return lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			divider,
			content,
			divider,
			footer,
		)
	}

	var footer string

	switch m.screen.(type) {
	case *DashboardModel:
		footer = "[/] Search   [↑/↓] Navigate   [Enter] Open   [Esc] Quit"

	case *ServerDetailModel:
		footer = "[H] History   [Esc] Back   [Q] Quit"
	}

	contentHeight := m.height - 5

	content := lipgloss.NewStyle().
		Width(m.width - 6).
		Height(contentHeight).
		Render(m.screen.View())

	footerView := lipgloss.NewStyle().
		Width(m.width-6).
		Height(2).
		Align(lipgloss.Center, lipgloss.Center).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		Render(helpStyle.Render(footer))

	inside := lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		footerView,
	)

	container := containerStyle.
		Width(m.width - 2).
		Height(m.height - 2)

	return container.Render(inside)
}
