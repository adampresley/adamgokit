package auth

import (
	"net/http"

	"github.com/adampresley/goth"
)

/*
GetFailureHandler wraps a SessionAuthHandler suitable for handling
Goth auth errors.
*/
func GetFailureHandler(handler AuthHandler) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		handler(w, r, nil, nil, goth.User{}, err)
	}
}
