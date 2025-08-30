package auth2

import "net/http"

type Authenticator[T any] interface {
	DestroySession(w http.ResponseWriter, r *http.Request) error
	Middleware(next http.Handler) http.Handler
	SaveSession(w http.ResponseWriter, r *http.Request, sessionValue T) error
}

type AuthHandler[T any] struct {
	authenticator Authenticator[T]
}

func New[T any](authenticator Authenticator[T]) *AuthHandler[T] {
	return &AuthHandler[T]{authenticator: authenticator}
}

func (a *AuthHandler[T]) DestroySession(w http.ResponseWriter, r *http.Request) error {
	return a.authenticator.DestroySession(w, r)
}

func (a *AuthHandler[T]) Middleware(next http.Handler) http.Handler {
	return a.authenticator.Middleware(next)
}

func (a *AuthHandler[T]) SaveSession(w http.ResponseWriter, r *http.Request, sessionValue T) error {
	return a.authenticator.SaveSession(w, r, sessionValue)
}
