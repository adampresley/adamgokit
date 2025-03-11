package sessions

import (
	gorillasessions "github.com/gorilla/sessions"
)

/*
NewCookieStore is a convenience method to initialize a cookie session store
with default options: MaxAge=0, Secure=false
*/
func NewCookieStore(sessionKey string) *gorillasessions.CookieStore {
	store := gorillasessions.NewCookieStore([]byte(sessionKey))
	return store
}

/*
NewCookieStoreWithOptions is a convenience method to initialize a cookie session store
with custom options.
*/
func NewCookieStoreWithOptions(sessionKey string, options *gorillasessions.Options) *gorillasessions.CookieStore {
	store := gorillasessions.NewCookieStore([]byte(sessionKey))
	store.Options = options
	return store
}
