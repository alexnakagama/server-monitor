package tui

import (
	"github.com/alexnakagama/server-monitor/internal/monitor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	username textinput.Model
	password textinput.Model
	client   *monitor.Client
}

func NewLoginModel(client *monitor.Client) LoginModel {
	username := textinput.New()
	username.Placeholder = "Username"
	username.Focus()

	password := textinput.New()
	password.Placeholder = "Password"
	password.EchoMode = textinput.EchoPassword

	return LoginModel{
		username: username,
		password: password,
		client:   client,
	}
}

func (m *LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.username.Focused() {
				m.username.Blur()
				m.password.Focus()
			} else {
				m.password.Blur()
				m.username.Focus()
			}

		case "enter":
			return m, nil
		}
	}

	var cmd tea.Cmd

	if m.username.Focused() {
		m.username, cmd = m.username.Update(msg)
	} else {
		m.password, cmd = m.password.Update(msg)
	}

	return m, cmd
}

func (m *LoginModel) View() string {
	return "Server Monitor\n\n" +
		"Username:\n" +
		m.username.View() +
		"\n\n" +
		"Password:\n" +
		m.password.View() +
		"\n\n" +
		"Press Enter to login"
}
