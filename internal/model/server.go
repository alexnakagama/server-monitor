package model

import "time"

type Server struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Hostname  string    `json:"hostname"`
	OS        string    `json:"os"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) Validate() error {}
