package auth

import (
	"context"
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type AgentAuthenticator interface {
	Authenticate(ctx context.Context, token string) (model.Agent, error)
}

func AgentAuthMiddleware(agentAuthenticator AgentAuthenticator, next http.Handler) http.Handler {}
