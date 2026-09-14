package main

import (
	"log"

	"github.com/alexnakagama/server-monitor/internal/monitor"
	"github.com/alexnakagama/server-monitor/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client := monitor.NewClient("http://localhost:8080")
	model := tui.NewLoginModel(client)

	p := tea.NewProgram(&model)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
