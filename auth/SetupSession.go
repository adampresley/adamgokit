package auth

import (
	"net/http"

	"github.com/adampresley/goth"
	"github.com/gorilla/sessions"
)

func SetupSession(w http.ResponseWriter, r *http.Request, config *AuthConfig, user goth.User) error {
	var (
		err     error
		session *sessions.Session
	)

	if session, err = config.Store.Get(r, config.SessionName); err != nil {
		return err
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

	if err = config.Store.Save(r, w, session); err != nil {
		return err
	}

	return nil
}

func DeleteSession(w http.ResponseWriter, r *http.Request, config *AuthConfig) error {
	var (
		err     error
		session *sessions.Session
	)

	if session, err = config.Store.Get(r, config.SessionName); err != nil {
		return err
	}

	session.Options.MaxAge = -1

	session.Values[EmailKey] = ""
	session.Values[FirstNameKey] = ""
	session.Values[LastNameKey] = ""
	session.Values[NameKey] = ""
	session.Values[ProviderKey] = ""
	session.Values[AvatarURLKey] = ""

	if err = config.Store.Save(r, w, session); err != nil {
		return err
	}

	return nil
}
