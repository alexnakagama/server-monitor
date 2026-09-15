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

var selectedStyle = lipgloss.NewStyle().
	Bold(true)
