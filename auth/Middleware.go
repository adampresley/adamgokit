package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/adampresley/adamgokit/sessions"
	"github.com/markbates/goth"
)

/*
This function returns a middleware that injects a Goth user from the session
into a context variable named "user".
*/
func getUserInjectorMiddleware(sessionConfig AuthConfig, middlewareConfig AuthMiddlewareConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				err        error
				userGetter sessions.SessionWrapper[goth.User]
				user       goth.User
			)

			/*
			 * Get the user from the session. If we don't have one, call
			 * the unauthorized handler. This allows callers to handle
			 * this how they want.
			 */
			userGetter = sessions.NewSessionWrapper[goth.User](sessionConfig.Store, sessionConfig.SessionName, UserSessionKey)

			if user, err = userGetter.Get(r); err != nil {
				slog.Error("user session not found", "error", err)
				middlewareConfig.UnauthorizedHandler(next).ServeHTTP(w, r)
				return
			}

			/*
			 * Put the user in the context for easy retrieval.
			 */
			ctx := context.WithValue(r.Context(), UserSessionKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
