package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/adampresley/goth"
	"github.com/adampresley/goth/gothic"
	"github.com/adampresley/goth/providers/apple"
	"github.com/adampresley/goth/providers/direct"
	"github.com/adampresley/goth/providers/facebook"
	"github.com/adampresley/goth/providers/google"
)

type Builder struct {
	Config AuthConfig
	Mux    *http.ServeMux
}

func NewBuilder(config AuthConfig, mux *http.ServeMux) *Builder {
	return &Builder{
		Config: config,
		Mux:    mux,
	}
}

func (b *Builder) WithApple(config OAuthConfig) *Builder {
	goth.UseProviders(
		apple.New(config.ClientID, config.ClientSecret, b.getCallbackURI("apple"), nil, config.Scopes...),
	)

	return b
}

func (b *Builder) WithDirect(config DirectConfig) *Builder {
	goth.UseProviders(
		direct.New(config.LoginURI, config.UserFetcher, config.CredChecker),
	)

	return b
}

func (b *Builder) WithFacebook(config OAuthConfig) *Builder {
	goth.UseProviders(
		facebook.New(config.ClientID, config.ClientSecret, b.getCallbackURI("facebook"), config.Scopes...),
	)

	return b
}

func (b *Builder) WithGoogle(config OAuthConfig) *Builder {
	goth.UseProviders(
		google.New(config.ClientID, config.ClientSecret, b.getCallbackURI("google"), config.Scopes...),
	)

	return b
}

func (b *Builder) Setup() *Builder {
	gothic.Store = b.Config.Store

	b.Mux.HandleFunc(
		fmt.Sprintf("GET %s/{provider}/callback", b.normalizeCallbackURIPrefix()),
		func(w http.ResponseWriter, r *http.Request) {
			successAuthCallback(w, r, &b.Config)
		},
	)

	b.Mux.HandleFunc(
		fmt.Sprintf("POST %s/{provider}/callback", b.normalizeCallbackURIPrefix()),
		func(w http.ResponseWriter, r *http.Request) {
			successAuthCallback(w, r, &b.Config)
		},
	)

	b.Mux.HandleFunc(fmt.Sprintf("GET %s/{provider}", b.normalizeCallbackURIPrefix()), func(w http.ResponseWriter, r *http.Request) {
		// This line is a workaround. The Goth library doesn't understand Go 1.22 path params
		// But it does understand having the provider in the context
		r = r.WithContext(context.WithValue(r.Context(), "provider", r.PathValue("provider")))
		gothic.BeginAuthHandler(w, r)
	})

	return b
}

func (b *Builder) normalizeCallbackURIPrefix() string {
	if strings.HasPrefix(b.Config.CallbackURIPrefix, "/") {
		return b.Config.CallbackURIPrefix
	}

	return "/" + b.Config.CallbackURIPrefix
}

func (b *Builder) getCallbackURI(provider string) string {
	redirectURI := fmt.Sprintf("%s%s/%s/callback", b.Config.BaseURL, b.normalizeCallbackURIPrefix(), provider)
	return redirectURI

}
