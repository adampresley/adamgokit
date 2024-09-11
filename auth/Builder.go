package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/adampresley/goth"
	"github.com/adampresley/goth/gothic"
	"github.com/adampresley/goth/providers/apple"
	"github.com/adampresley/goth/providers/direct"
	"github.com/adampresley/goth/providers/facebook"
	"github.com/adampresley/goth/providers/google"
	"github.com/gorilla/sessions"
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

	callbackHandler := func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			user    goth.User
			session *sessions.Session
		)

		user, err = gothic.CompleteUserAuth(w, r)

		if err != nil {
			slog.Error("user authorization failed due to an error", "error", err)
			http.Redirect(w, r, b.Config.ErrorPath+"?message="+err.Error(), http.StatusFound)
			return
		}

		if session, err = b.Config.Store.Get(r, b.Config.SessionName); err != nil {
			slog.Error("could not get session", "error", err)
			http.Redirect(w, r, b.Config.ErrorPath+"?message="+err.Error(), http.StatusFound)
			return
		}

		/*
		 * Store important information in the session.
		 */
		session.Values[EmailKey] = user.Email
		session.Values[FirstNameKey] = user.FirstName
		session.Values[LastNameKey] = user.LastName
		session.Values[NameKey] = user.Name
		session.Values[ProviderKey] = user.Provider
		session.Values[AvatarURLKey] = user.AvatarURL

		if err = b.Config.Store.Save(r, w, session); err != nil {
			slog.Error("could not save user in session", "error", err)
			http.Redirect(w, r, b.Config.ErrorPath+"?message="+err.Error(), http.StatusFound)
			return
		}

		b.Config.Handler(w, r, b.Config.Store, session, user, nil)
	}

	b.Mux.HandleFunc(
		fmt.Sprintf("GET %s/{provider}/callback", b.normalizeCallbackURIPrefix()),
		callbackHandler,
	)

	b.Mux.HandleFunc(
		fmt.Sprintf("POST %s/{provider}/callback", b.normalizeCallbackURIPrefix()),
		callbackHandler,
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
