package auth2_test

import (
	"net/http"
	"testing"

	"github.com/adampresley/adamgokit/auth2"
	"github.com/stretchr/testify/assert"
)

func TestWithContextKey(t *testing.T) {
	want := "test-context-key"
	options := &auth2.AuthenticatorOptions{}

	option := auth2.WithContextKey(want)
	option(options)

	assert.Equal(t, want, options.ContextKey)
}

func TestWithDebug(t *testing.T) {
	options := &auth2.AuthenticatorOptions{}

	option := auth2.WithDebug(true)
	option(options)

	assert.True(t, options.Debug)

	option = auth2.WithDebug(false)
	option(options)

	assert.False(t, options.Debug)
}

func TestWithErrorFunc(t *testing.T) {
	options := &auth2.AuthenticatorOptions{}
	errorFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	option := auth2.WithErrorFunc(errorFunc)
	option(options)

	assert.NotNil(t, options.ErrorFunc)
}

func TestWithExcludedPaths(t *testing.T) {
	want := []string{"/health", "/metrics", "/ping"}
	options := &auth2.AuthenticatorOptions{}

	option := auth2.WithExcludedPaths(want)
	option(options)

	assert.Equal(t, want, options.ExcludedPaths)
}

func TestWithExcludedPathsExact(t *testing.T) {
	options := &auth2.AuthenticatorOptions{}

	option := auth2.WithExcludedPathsExact(true)
	option(options)

	assert.True(t, options.ExcludedPathsExact)

	option = auth2.WithExcludedPathsExact(false)
	option(options)

	assert.False(t, options.ExcludedPathsExact)
}

func TestWithRedirectURL(t *testing.T) {
	want := "https://example.com/login"
	options := &auth2.AuthenticatorOptions{}

	option := auth2.WithRedirectURL(want)
	option(options)

	assert.Equal(t, want, options.RedirectURL)
}

func TestWithResponder(t *testing.T) {
	options := &auth2.AuthenticatorOptions{}
	responderFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusUnauthorized)
	}

	option := auth2.WithResponder(responderFunc)
	option(options)

	assert.NotNil(t, options.ResponderFunc)
}

