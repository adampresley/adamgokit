package auth

import (
	"net/http"

	"github.com/markbates/goth"
)

/*
GetFailureHandler wraps a SessionAuthHandler suitable for handling
Goth auth errors.
*/
func GetFailureHandler(handler SessionAuthHandler) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		handler(w, r, nil, nil, goth.User{}, err)
	}
}
