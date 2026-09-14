package tui

import tea "github.com/charmbracelet/bubbletea"

type DashboardModel struct{}

func NewDashBoardModel() DashboardModel {
	return DashboardModel{}
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func quit() tea.Cmd {
	return func() tea.Msg {
		return tea.Quit()
	}
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, quit()
		}
	}

	return m, nil
}

func (m DashboardModel) View() string {
	return "Server Monitor\n\n" +
		"Dashboard\n\n" +
		"Logged in successfully.\n\n" +
		"Press q to quit\n"
}
