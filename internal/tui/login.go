package tui

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/monitor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	username textinput.Model
	password textinput.Model
	client   *monitor.Client
	err      error
	loading  bool
	width    int
	height   int
}

type loginResultMessage struct {
	err error
}

func (m *LoginModel) login() tea.Cmd {
	return func() tea.Msg {
		err := m.client.Login(
			context.Background(),
			m.username.Value(),
			m.password.Value(),
		)
		if err != nil {
			return loginResultMessage{err: err}
		}

		return loginResultMessage{}
	}
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil

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
			m.loading = true
			m.err = nil

			return m, m.login()
		}

	case loginResultMessage:
		m.loading = false

		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		return m, func() tea.Msg {
			return dashboardMessage{
				client: m.client,
			}
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
	view := ""

	view += titleStyle.Render("LOGIN") + "\n\n"

	view += labelStyle.Render("Username") + "\n"
	view += m.username.View() + "\n\n"

	view += labelStyle.Render("Password") + "\n"
	view += m.password.View() + "\n\n"

	if m.loading {
		view += helpStyle.Render("Logging in...") + "\n"
	} else if m.err != nil {
		view += errorStyle.Render("Login failed: "+m.err.Error()) + "\n"
	} else {
		view += helpStyle.Render("Press Enter to login") + "\n"
	}

	return loginBoxStyle.Render(view)
}
