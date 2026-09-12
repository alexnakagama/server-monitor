package model

import (
	"errors"
	"strings"
	"testing"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

func TestValidateAgentName(t *testing.T) {
	tests := []struct {
		name    string
		agent   string
		wantErr error
	}{
		{
			name:    "empty agent name is required",
			agent:   "",
			wantErr: errors_custom.ErrNameRequired,
		},
		{
			name:    "agent name of 100 characters is valid",
			agent:   strings.Repeat("a", 100),
			wantErr: nil,
		},
		{
			name:    "agent name longer than 100 characters is too long",
			agent:   strings.Repeat("a", 101),
			wantErr: errors_custom.ErrUsernameTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateAgentName(tt.agent)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("ValidateAgentName(%q) error = %v, want %v", tt.agent, gotErr, tt.wantErr)
			}
		})
	}
}

func TestValidateServerID(t *testing.T) {
	tests := []struct {
		name     string
		serverID int
		wantErr  error
	}{
		{
			name:     "zero server id is invalid",
			serverID: 0,
			wantErr:  errors_custom.ErrInvalidServerID,
		},
		{
			name:     "negative server id is invalid",
			serverID: -1,
			wantErr:  errors_custom.ErrInvalidServerID,
		},
		{
			name:     "positive server id is valid",
			serverID: 1,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateServerID(tt.serverID)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("ValidateServerID(%d) error = %v, want %v", tt.serverID, gotErr, tt.wantErr)
			}
		})
	}
}

func TestAgentValidate(t *testing.T) {
	tests := []struct {
		name    string
		agent   *Agent
		wantErr error
	}{
		{
			name: "valid agent passes validation",
			agent: &Agent{
				Name:     "agent-1",
				ServerID: 1,
			},
			wantErr: nil,
		},
		{
			name: "agent with empty name fails validation",
			agent: &Agent{
				Name:     "",
				ServerID: 1,
			},
			wantErr: errors_custom.ErrNameRequired,
		},
		{
			name: "agent with invalid server id fails validation",
			agent: &Agent{
				Name:     "agent-1",
				ServerID: 0,
			},
			wantErr: errors_custom.ErrInvalidServerID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.agent.Validate()

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Agent.Validate() error = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
