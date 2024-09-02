package auth

import (
	"net/http"

	"github.com/markbates/goth"
)

/*
GetUserFromContext gets the Goth user from context and
returns it as a concrete type. If it cannot be cast,
and empty User struct is returned.
*/
func GetUserFromContext(r *http.Request) goth.User {
	if result, ok := r.Context().Value(UserSessionKey).(goth.User); ok {
		return result
	}

	return goth.User{}
}
