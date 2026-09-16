package tui

import "github.com/charmbracelet/lipgloss"

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 1)

var sectionStyle = lipgloss.NewStyle().
	Bold(true)

var labelStyle = lipgloss.NewStyle().
	Width(12)

var valueStyle = lipgloss.NewStyle()

var helpStyle = lipgloss.NewStyle().
	Faint(true)

var containerStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 2)

var headerStyle = lipgloss.NewStyle().
	Bold(true)

var onlineStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10"))

var offlineStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9"))

var nameColumnStyle = lipgloss.NewStyle().
	Width(20)

var osColumnStyle = lipgloss.NewStyle().
	Width(10)

var prefixStyle = lipgloss.NewStyle().
	Width(2)

var sectionBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(0, 2).
	MarginBottom(1)
