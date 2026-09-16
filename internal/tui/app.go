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
