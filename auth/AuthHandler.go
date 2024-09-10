package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/adampresley/goth"
)

type AuthHandler func(w http.ResponseWriter, r *http.Request, store sessions.Store, session *sessions.Session, user goth.User, err error)
