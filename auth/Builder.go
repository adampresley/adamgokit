package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/adampresley/adamgokit/httphelpers"
	"github.com/gorilla/sessions"
	"github.com/adampresley/goth"
	"github.com/adampresley/goth/gothic"
	"github.com/adampresley/goth/providers/apple"
	"github.com/adampresley/goth/providers/facebook"
	"github.com/adampresley/goth/providers/google"
)

type Builder struct {
	Config AuthConfig
	Mux    *http.ServeMux

	// Since Goth doesn't support "direct" auth schemes yet,
	// we'll do it ourselves.
	useDirect    bool
	directConfig DirectConfig
}

func NewBuilder(config AuthConfig, mux *http.ServeMux) *Builder {
	return &Builder{
		Config: config,
		Mux:    mux,

		useDirect: false,
	}
}

func (b *Builder) WithApple(config OAuthConfig) *Builder {
	goth.UseProviders(
		apple.New(config.ClientID, config.ClientSecret, b.getCallbackURI("apple"), nil, config.Scopes...),
	)

	return b
}

func (b *Builder) WithDirect(config DirectConfig) *Builder {
	b.useDirect = true
	b.directConfig = config
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
			http.Redirect(w, r, b.Config.ErrorPath, http.StatusTemporaryRedirect)
			return
		}

		b.Config.Handler(w, r, b.Config.Store, session, user, nil)
	})

	b.Mux.HandleFunc(fmt.Sprintf("GET %s/{provider}", b.normalizeCallbackURIPrefix()), func(w http.ResponseWriter, r *http.Request) {
		// This line is a workaround. The Goth library doesn't understand Go 1.22 path params
		r = r.WithContext(context.WithValue(r.Context(), "provider", r.PathValue("provider")))
		gothic.BeginAuthHandler(w, r)
	})

	if b.useDirect {
		b.Mux.HandleFunc(fmt.Sprintf("POST %s/direct", b.normalizeCallbackURIPrefix()), func(w http.ResponseWriter, r *http.Request) {
			var (
				err     error
				user    goth.User
				session *sessions.Session
			)

			orgID := httphelpers.GetFromRequest[string](r, "orgID")

			loginInput := DirectUserLoginInput{
				UserName: httphelpers.GetFromRequest[string](r, "userName"),
				Password: httphelpers.GetFromRequest[string](r, "password"),
				OrgID:    orgID,
			}

			user, err = b.directConfig.UserValidator(loginInput)

			if err != nil {
				GetFailureHandler(b.Config.Handler)(w, r, err)
				return
			}

			if session, err = b.Config.Store.Get(r, b.Config.SessionName); err != nil {
				slog.Error("could not get session", "error", err)
				http.Redirect(w, r, b.Config.ErrorPath, http.StatusTemporaryRedirect)
				return
			}

			/*
			 * Store important information in the session.
			 */
			session.Values[EmailKey] = user.Email
			session.Values[FirstNameKey] = user.FirstName
			session.Values[LastNameKey] = user.LastName
			session.Values[NameKey] = user.Name
			session.Values[ProviderKey] = "direct"
			session.Values[AvatarURLKey] = user.AvatarURL
			session.Values[OrgIDKey] = orgID

			if err = b.Config.Store.Save(r, w, session); err != nil {
				slog.Error("could not save user in session", "error", err)
				http.Redirect(w, r, b.Config.ErrorPath, http.StatusTemporaryRedirect)
				return
			}

			b.Config.Handler(w, r, b.Config.Store, session, user, nil)
		})
	}

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
