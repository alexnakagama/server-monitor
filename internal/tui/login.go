package tui

import "github.com/charmbracelet/bubbles/textinput"

type LoginModel struct {
	username textinput.Model
	password textinput.Model
}
