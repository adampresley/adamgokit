package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
)

type Builder struct {
	Config SessionAuthConfig
	Mux    *http.ServeMux
}

func NewBuilder(config SessionAuthConfig, mux *http.ServeMux) *Builder {
	return &Builder{
		Config: config,
		Mux:    mux,
	}
}

func (b *Builder) WithApple(scopes []string) *Builder {
	redirectURI := fmt.Sprintf("%s/facebook/callback", b.Config.CallbackURIPrefix)

	goth.UseProviders(
		apple.New(b.Config.ClientKey, b.Config.ClientSecret, redirectURI, nil, scopes...),
	)

	return b
}

func (b *Builder) WithFacebook() *Builder {
	redirectURI := fmt.Sprintf("%s/facebook/callback", b.Config.CallbackURIPrefix)

	goth.UseProviders(
		facebook.New(b.Config.ClientKey, b.Config.ClientSecret, redirectURI),
	)

	return b
}

func (b *Builder) WithGoogle(scopes []string) *Builder {
	redirectURI := fmt.Sprintf("%s/google/callback", b.Config.CallbackURIPrefix)

	goth.UseProviders(
		google.New(b.Config.ClientKey, b.Config.ClientSecret, redirectURI, scopes...),
	)

	return b
}

func (b *Builder) Setup() *Builder {
	gothic.Store = b.Config.Store

	b.Mux.HandleFunc(fmt.Sprintf("%s/{provider}/callback", b.Config.CallbackURIPrefix), func(w http.ResponseWriter, r *http.Request) {
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

		b.Config.Handler(w, r, b.Config.Store, session, user, nil)
	})

	b.Mux.HandleFunc(fmt.Sprintf("%s/{provider}", b.Config.CallbackURIPrefix), func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	return b
}
