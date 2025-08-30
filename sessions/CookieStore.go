package sessions

import (
	"net/http"
	"time"

	gorillasessions "github.com/gorilla/sessions"
)

type CookieStoreOption func(options *gorillasessions.Options)

/*
NewCookieStore is a convenience method to initialize a cookie session store
with custom options.
*/
func NewCookieStore(sessionKey string, options ...CookieStoreOption) *gorillasessions.CookieStore {
	opts := &gorillasessions.Options{
		SameSite: http.SameSiteDefaultMode,
	}

	for _, opt := range options {
		opt(opts)
	}

	store := gorillasessions.NewCookieStore([]byte(sessionKey))
	store.Options = opts
	return store
}

func WithSecure(secure bool) CookieStoreOption {
	return func(options *gorillasessions.Options) {
		options.Secure = secure
	}
}

func WithMaxAge(maxAge time.Duration) CookieStoreOption {
	return func(options *gorillasessions.Options) {
		options.MaxAge = int(maxAge.Seconds())
	}
}

func WithSameSite(sameSite http.SameSite) CookieStoreOption {
	return func(options *gorillasessions.Options) {
		options.SameSite = sameSite
	}
}

func WithDomain(domain string) CookieStoreOption {
	return func(options *gorillasessions.Options) {
		options.Domain = domain
	}
}

func WithHttpOnly(httpOnly bool) CookieStoreOption {
	return func(options *gorillasessions.Options) {
		options.HttpOnly = httpOnly
	}
}
