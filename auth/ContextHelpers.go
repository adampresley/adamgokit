package auth

import (
	"net/http"
)

/*
GetEmailFromContext gets the Goth user email from context and
returns it as a concrete type.
*/
func GetEmailFromContext(r *http.Request) string {
	if result, ok := r.Context().Value(EmailKey).(string); ok {
		return result
	}

	return ""
}
