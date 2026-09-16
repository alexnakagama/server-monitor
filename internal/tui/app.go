package tui

import tea "github.com/charmbracelet/bubbletea"

type AppModel struct {
	width  int
	height int
	screen tea.Model
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
	}

	var cmd tea.Cmd

	m.screen, cmd = m.screen.Update(msg)

	return m, cmd
}

func (m *AppModel) View() string {
	return m.screen.View()
}
