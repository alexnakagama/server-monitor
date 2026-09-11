package auth

import "net/http"

func RequireRole(role string, next http.Handler) http.Handler {}
