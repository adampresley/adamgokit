/*
auth package provides wrappers for authentication.

Example:

	  storeName := "example"
	  store := gorillasessions.NewCookieStore([]byte("secret"))
	  userSession := session.NewSessionWrapper[goth.User](store, storeName)

		  authSuccessHandler := func(w http.ResponseWriter, r *http.Request, store sessions.Store, session *sessions.Session, user goth.User, err error) {
		    if err != nil {
		      httphelpers.JsonErrorMessage(r, http.StatusInternalServerError, "can't do it")
		      return
		    }

		    // All is good. Set up our session
		    userSession.Set(r, user)

		    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		  }

		  getAuthMiddleware := func(s session.Session, config auth.SessionAuthConfig) func(http.Handler) http.Handler {
		    return func(next http.Handler) http.Handler {
		      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		        // Get our user from session and put it in context. Here we can do other stuff too.
		        user, err := s.Get(r)

		        // Here we can check err, or check excluded paths from config, etc...

		        ctx := context.WithValue(r.Context(), "user", user)
		        next.ServeHTTP(w, r.WithContext(ctx))
		      })
		    }
		  }

		  config := auth.SessionAuthConfig{
		    SessionAuthConfig: auth.SessionAuthConfig{
		      Handler: authSuccessHandler,
		      ErrorPath: "/error",
		      ExcludedPaths: []string{"/login", "/logout"},
		      SessionName: "example",
		    },
		    CallbackURIPrefix: "/auth",
		    ClientID: "",
		    ClientSecret: "",
		  }

	  auth.NewBuilder(config, mux).
		  WithApple(auth.DefaultAppleScopes).
		  WithGoogle(auth.DefaultGoogleScopes).
		  WithFacebook().
		  Setup()

	  authMiddleware := getAuthMiddleware(userSession, config)

	  mux.HandleFunc("GET /test", authMiddleware(myHandler))
*/
package auth
