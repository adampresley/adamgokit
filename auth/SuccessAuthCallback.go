package auth

import (
	"log/slog"
	"net/http"

	"github.com/adampresley/goth"
	"github.com/adampresley/goth/gothic"
	"github.com/gorilla/sessions"
)

func successAuthCallback(w http.ResponseWriter, r *http.Request, config *AuthConfig) {
	var (
		err     error
		user    goth.User
		session *sessions.Session
	)

	if session, err = config.Store.Get(r, config.SessionName); err != nil {
		slog.Error("could not get session", "error", err)
		http.Redirect(w, r, config.ErrorPath+"?message="+err.Error(), http.StatusFound)
		return
	}

	user, err = gothic.CompleteUserAuth(w, r)

	if err != nil {
		slog.Error("user authorization failed due to an error", "error", err)
		config.Handler(w, r, config.Store, session, user, err)
		return
	}

	if err = SetupSession(w, r, config, user); err != nil {
		slog.Error("could not save user in session", "error", err)
		http.Redirect(w, r, config.ErrorPath+"?message="+err.Error(), http.StatusFound)
		return
	}

	config.Handler(w, r, config.Store, session, user, nil)
}
