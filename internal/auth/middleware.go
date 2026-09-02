package auth

import "context"

type contextKey string

const userIDKey contextKey = "userID"

func UserIDFromContext(ctx context.Context) (int, bool) {}
