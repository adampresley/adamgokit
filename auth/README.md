# Authentication Tools

This package provides wrappers for handling various authentication methods. It
makes use of the following third-party libraries:

- [Gorilla Sessions](https://github.com/gorilla/sessions)
- [Goth](https://github.com/markbates/goth)

## Google

```go
func AuthCallbackHandler(w http.ResponseWriter, r *http.Request, store gorillasessions.Store, session *gorillasessions.Session, user goth.User, err error) {
  var (
    ident  *identity.Identity
    dbUser *identity.User
  )

  /*
   * First, if there was an error, do something about it.
   */
  if err != nil {
    // Here you should log, or redirect, or return JSON
    return
  }

  /*
   * This is where you might do something like create identity/user records
   * if the user is logging in through OAuth.
   */

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

/*
 * Setup session store
 */
store, storeCleaner, err := sessions.NewPGStore(config.DSN, config.SessionKey)

if err != nil {
  slog.Error("error setting up session storage. aborting", "error", err)
  os.Exit(-1)
}

defer storeCleaner()

/*
 * Setup auth components
 */
authConfig := auth.AuthConfig{
  BaseURL:           config.BaseURL,
  CallbackURIPrefix: "/auth",
  Handler:           AuthCallbackHandler,
  ErrorPath:         "/error",
  SessionName:       sessionName,
  Store:             store,
}

// Here, "m" is a http.ServeMux
auth.NewBuilder(authConfig, m).
  WithGoogle(auth.OAuthConfig{
    ClientID:     config.GoogleClientID,
    ClientSecret: config.GoogleClientSecret,
    Scopes:       auth.DefaultGoogleScopes,
  }).
  WithFacebook(auth.OAuthConfig{
    ClientID:     config.FacebookClientID,
    ClientSecret: config.FacebookClientSecret,
    Scopes:       auth.DefaultFacebookScopes,
  }).
  Setup()
```
