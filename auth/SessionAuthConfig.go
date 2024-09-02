package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
)

type SessionAuthHandler func(w http.ResponseWriter, r *http.Request, store sessions.Store, session *sessions.Session, user goth.User, err error)

/*
SessionAuthConfig provides configuration information for authentication.
It contains settings for different types of authentication.
*/
type SessionAuthConfig struct {
	CallbackURIPrefix string
	ClientKey         string
	ClientSecret      string
	Handler           SessionAuthHandler
	ErrorPath         string
	ExcludedPaths     []string
	SessionName       string
	Store             sessions.Store
}
