package auth

import (
	"github.com/gorilla/sessions"
)

/*
SessionAuthConfig provides configuration information for authentication.
It contains settings for different types of authentication.
*/
type AuthConfig struct {
	BaseURL           string
	CallbackURIPrefix string
	Debug             bool
	Handler           AuthHandler
	ErrorPath         string
	SessionName       string
	Store             sessions.Store
}
