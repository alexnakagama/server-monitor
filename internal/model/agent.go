package model

import "time"

type Agent struct {
	ID        int       `json:"id"`
	ServerID  int       `json:"server_id"`
	Name      string    `json:"name"`
	TokenHash string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Agent) Validate() error {}
