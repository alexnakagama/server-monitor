package main

import (
	"log"

	"github.com/alexnakagama/server-monitor/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := tui.NewLoginModel()

	p := tea.NewProgram(&model)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
