package tui

import "github.com/charmbracelet/bubbles/textinput"

type LoginModel struct {
	username textinput.Model
	password textinput.Model
}

func NewLoginModel() LoginModel {
	username := textinput.New()
	username.Placeholder = "username"
	username.Focus()

	password := textinput.New()
	password.Placeholder = "username"
	password.Focus()

	return LoginModel{
		username: username,
		password: password,
	}
}
