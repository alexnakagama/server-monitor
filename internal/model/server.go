package model

import (
	"errors"
	"time"
)

type Server struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Hostname  string    `json:"hostname"`
	OS        string    `json:"os"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) Validate() error {
	if s.Name == "" {
		return errors.New("name is required")
	}

	if s.Hostname == "" {
		return errors.New("hostname is required")
	}

	if s.OS == "" {
		return errors.New("os is required")
	}

	return nil
}

func (s *Server) ValidateName() error {}
