package auth

import "net/http"

type AuthMiddlewareConfig struct {
	ExcludedPaths       []string
	UnauthorizedHandler func(http.Handler) http.Handler
}
