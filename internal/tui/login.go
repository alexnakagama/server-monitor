package tui

import "github.com/charmbracelet/bubbles/textinput"

type LoginModel struct {
	username textinput.Model
	password textinput.Model
}

func NewLoginModel() LoginModel {
	username := textinput.New()
	username.Placeholder = "Username"
	username.Focus()

	password := textinput.New()
	password.Placeholder = "Password"
	password.EchoMode = textinput.EchoPassword

	return LoginModel{
		username: username,
		password: password,
	}
}
