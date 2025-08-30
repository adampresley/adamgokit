package auth2

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	gorillasessions "github.com/gorilla/sessions"
)

type UserNameAndPasswordProvider[T any] struct {
	authenticatorOptions *AuthenticatorOptions
	contextKey           string
	debug                bool
	sessionStore         gorillasessions.Store
	sessionName          string
	sessionKey           string
}

func UserNameAndPassword[T any](sessionStore gorillasessions.Store, sessionName, sessionKey string, options ...AuthenticatorOption) *UserNameAndPasswordProvider[T] {
	result := &UserNameAndPasswordProvider[T]{
		authenticatorOptions: &AuthenticatorOptions{},
		contextKey:           "session",
		debug:                false,
		sessionStore:         sessionStore,
		sessionName:          sessionName,
		sessionKey:           sessionKey,
	}

	for _, option := range options {
		option(result.authenticatorOptions)
	}

	if result.authenticatorOptions.ContextKey != "" {
		result.contextKey = result.authenticatorOptions.ContextKey
	}

	if result.authenticatorOptions.Debug {
		result.debug = result.authenticatorOptions.Debug
	}

	return result
}

func (a *UserNameAndPasswordProvider[T]) DestroySession(w http.ResponseWriter, r *http.Request) error {
	var (
		err     error
		session *gorillasessions.Session
	)

	if session, err = a.sessionStore.Get(r, a.sessionName); err != nil {
		return fmt.Errorf("error getting session: %w", err)
	}

	session.Options.MaxAge = -1
	delete(session.Values, a.sessionKey)

	return session.Save(r, w)
}

func (a *UserNameAndPasswordProvider[T]) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			session *gorillasessions.Session
		)

		a.debugf("starting authentication middleware", "path", r.URL.Path)

		/*
		 * If excluded paths have been provided, check if the current request path matches any of them.
		 */
		if len(a.authenticatorOptions.ExcludedPaths) > 0 {
			if a.authenticatorOptions.ExcludedPathsExact {
				if a.pathExactlyMatchesExcludedPath(r.URL.Path) {
					next.ServeHTTP(w, r)
					return
				}
			} else {
				if a.pathContainsExcludedPath(r.URL.Path) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		/*
		 * Get the session from the request context. If we don't have a session,
		 * determine what we should do based on the provided options.
		 */
		if session, err = a.sessionStore.Get(r, a.sessionName); err != nil {
			if a.authenticatorOptions.ErrorFunc != nil {
				a.authenticatorOptions.ErrorFunc(w, r, err)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if _, ok := session.Values[a.sessionKey]; !ok {
			if a.authenticatorOptions.RedirectURL != "" {
				http.Redirect(w, r, a.authenticatorOptions.RedirectURL, http.StatusSeeOther)
				return
			}

			if a.authenticatorOptions.ResponderFunc != nil {
				a.authenticatorOptions.ResponderFunc(w, r, ErrSessionKeyNotFound)
				return
			}

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rawSessionValue := session.Values[a.sessionKey]
		sessionValue, ok := any(rawSessionValue).(T)

		if !ok {
			a.debugf("error converting session value to expected type", "sessionValue", rawSessionValue)

			if a.authenticatorOptions.ErrorFunc != nil {
				a.authenticatorOptions.ErrorFunc(w, r, ErrSessionValueNotConvertible)
				return
			}

			http.Error(w, "Error converting session value to expected type", http.StatusInternalServerError)
			return
		}

		/*
		 * Add the session value to the request context
		 */
		a.debugf("added session value to request context", "sessionValue", sessionValue)
		ctx := context.WithValue(r.Context(), a.contextKey, sessionValue)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *UserNameAndPasswordProvider[T]) SaveSession(w http.ResponseWriter, r *http.Request, sessionValue T) error {
	var (
		err     error
		session *gorillasessions.Session
	)

	if session, err = a.sessionStore.Get(r, a.sessionName); err != nil {
		return fmt.Errorf("error getting session: %w", err)
	}

	session.Values[a.sessionKey] = sessionValue
	return session.Save(r, w)
}

func (a *UserNameAndPasswordProvider[T]) pathContainsExcludedPath(path string) bool {
	for _, excludedPath := range a.authenticatorOptions.ExcludedPaths {
		if strings.HasPrefix(path, excludedPath) {
			a.debugf("path contains excluded path", "path", path)
			return true
		}
	}

	return false
}

func (a *UserNameAndPasswordProvider[T]) pathExactlyMatchesExcludedPath(path string) bool {
	if slices.Contains(a.authenticatorOptions.ExcludedPaths, path) {
		a.debugf("path matches excluded path exactly", "path", path)
		return true
	}

	return false
}

func (a *UserNameAndPasswordProvider[T]) debugf(message string, args ...any) {
	if a.debug {
		slog.Info(message, args...)
	}
}
