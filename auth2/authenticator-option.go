package auth2

import "net/http"

type AuthenticatorOption func(options *AuthenticatorOptions)

type AuthenticatorOptions struct {
	ContextKey         string
	Debug              bool
	ErrorFunc          func(w http.ResponseWriter, r *http.Request, err error)
	ExcludedPaths      []string
	ExcludedPathsExact bool
	RedirectURL        string
	ResponderFunc      func(w http.ResponseWriter, r *http.Request, err error)
}

func WithContextKey(contextKey string) AuthenticatorOption {
	return func(options *AuthenticatorOptions) {
		options.ContextKey = contextKey
	}
}

func WithDebug(debug bool) AuthenticatorOption {
	return func(options *AuthenticatorOptions) {
		options.Debug = debug
	}
}

func WithErrorFunc(errorFunc func(w http.ResponseWriter, r *http.Request, err error)) AuthenticatorOption {
	return func(options *AuthenticatorOptions) {
		options.ErrorFunc = errorFunc
	}
}

func WithExcludedPaths(excludedPaths []string) AuthenticatorOption {
	return func(options *AuthenticatorOptions) {
		options.ExcludedPaths = excludedPaths
	}
}

func WithExcludedPathsExact(excludedPathsExact bool) AuthenticatorOption {
	return func(options *AuthenticatorOptions) {
		options.ExcludedPathsExact = excludedPathsExact
	}
}

func WithRedirectURL(redirectURL string) AuthenticatorOption {
	return func(options *AuthenticatorOptions) {
		options.RedirectURL = redirectURL
	}
}

func WithResponder(responderFunc func(w http.ResponseWriter, r *http.Request, err error)) AuthenticatorOption {
	return func(options *AuthenticatorOptions) {
		options.ResponderFunc = responderFunc
	}
}
