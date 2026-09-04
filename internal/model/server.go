package model

import (
	"errors"
	"strings"
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

func (s *Server) ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	return nil
}

func (s *Server) ValidateHostname(hostname string) error {
	if strings.TrimSpace(hostname) == "" {
		return errors.New("hostname is required")
	}

	return nil
}

func (s *Server) ValidateOS(os string) error {
	if strings.TrimSpace(os) == "" {
		return errors.New("os is required")
	}

	return nil
}
