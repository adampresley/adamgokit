package auth

import (
	"net/http"
	"strings"
)

type AuthMiddlewareConfig struct {
	ExcludedPaths []string
	Middleware    func(http.Handler) http.Handler
}

func (am AuthMiddlewareConfig) IsPathExcluded(path string, exactMatch bool) bool {
	excluded := false

	for _, ep := range am.ExcludedPaths {
		if exactMatch && ep == path {
			excluded = true
			break
		} else if !exactMatch && strings.HasPrefix(path, ep) {
			excluded = true
			break
		}
	}

	return excluded
}
