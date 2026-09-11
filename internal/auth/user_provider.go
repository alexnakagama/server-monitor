package auth

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type UserProvider interface {
	GetByID(ctx context.Context, id int) (model.User, error)
}
