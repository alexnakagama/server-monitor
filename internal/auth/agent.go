package auth

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type AgentAuthenticator interface {
	Authenticate(ctx context.Context, token string) (model.Agent, error)
}
