package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
)

type AuthHandler func(w http.ResponseWriter, r *http.Request, store sessions.Store, session *sessions.Session, user goth.User, err error)

/*
SessionAuthConfig provides configuration information for authentication.
It contains settings for different types of authentication.
*/
type AuthConfig struct {
	CallbackURIPrefix string
	Handler           AuthHandler
	ErrorPath         string
	SessionName       string
	Store             sessions.Store
}
