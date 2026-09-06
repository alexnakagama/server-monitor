package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type AgentAuthenticator interface {
	Authenticate(ctx context.Context, token string) (model.Agent, error)
}

type agentIDKey string
type serverIDKey string

const (
	agentIDContextKey  agentIDKey  = "agentID"
	serverIDContextKey serverIDKey = "serverID"
)

func AgentIDFromContext(ctx context.Context) (int, bool) {}

func AgentAuthMiddleware(agentAuthenticator AgentAuthenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHandler := r.Header.Get("Authorization")

		if authHandler == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHandler, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authotization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		agent, err := agentAuthenticator.Authenticate(r.Context(), token)
		if err != nil {
			http.Error(w, "invalid agent token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), agentIDContextKey, agent.ID)
		ctx = context.WithValue(r.Context(), serverIDContextKey, agent.ServerID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
