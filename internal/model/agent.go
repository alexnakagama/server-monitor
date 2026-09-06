package model

import (
	"time"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

type Agent struct {
	ID        int       `json:"id"`
	ServerID  int       `json:"server_id"`
	Name      string    `json:"name"`
	TokenHash string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Agent) Validate() error {}

func ValidateAgentName(name string) error {
	if name == "" {
		return errors_custom.ErrNameRequired
	}

	return nil
}
