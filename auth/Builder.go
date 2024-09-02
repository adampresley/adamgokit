package auth

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
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
	redirectURI := fmt.Sprintf("%s/facebook/callback", b.normalizeCallbackURIPrefix())

	goth.UseProviders(
		apple.New(config.ClientID, config.ClientSecret, redirectURI, nil, config.Scopes...),
	)

	return b
}

func (b *Builder) WithFacebook(config OAuthConfig) *Builder {
	redirectURI := fmt.Sprintf("%s/facebook/callback", b.normalizeCallbackURIPrefix())

	goth.UseProviders(
		facebook.New(config.ClientID, config.ClientSecret, redirectURI, config.Scopes...),
	)

	return b
}

func (b *Builder) WithGoogle(config OAuthConfig) *Builder {
	redirectURI := fmt.Sprintf("%s/google/callback", b.normalizeCallbackURIPrefix())

	goth.UseProviders(
		google.New(config.ClientID, config.ClientSecret, redirectURI, config.Scopes...),
	)

	return b
}

func (b *Builder) Setup() *Builder {
	gothic.Store = b.Config.Store

	b.Mux.HandleFunc(fmt.Sprintf("GET %s/{provider}/callback", b.normalizeCallbackURIPrefix()), func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			user    goth.User
			session *sessions.Session
		)

		user, err = gothic.CompleteUserAuth(w, r)

		if err != nil {
			GetFailureHandler(b.Config.Handler)(w, r, err)
			return
		}

		if session, err = b.Config.Store.Get(r, b.Config.SessionName); err != nil {
			slog.Error("could not get session", "error", err)
			http.Redirect(w, r, b.Config.ErrorPath, http.StatusTemporaryRedirect)
			return
		}

		session.Values[UserSessionKey] = user

		if err = b.Config.Store.Save(r, w, session); err != nil {
			slog.Error("could not save user in session", "error", err)
			http.Redirect(w, r, b.Config.ErrorPath, http.StatusTemporaryRedirect)
			return
		}

		b.Config.Handler(w, r, b.Config.Store, session, user, nil)
	})

	b.Mux.HandleFunc(fmt.Sprintf("GET %s/{provider}", b.normalizeCallbackURIPrefix()), func(w http.ResponseWriter, r *http.Request) {
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
