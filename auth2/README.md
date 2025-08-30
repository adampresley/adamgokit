# Authentication Tools

This package provides wrappers for handling various authentication methods. It
makes use of the following third-party libraries:

- [Gorilla Sessions](https://github.com/gorilla/sessions)

## Key Concepts

This authentication library is centered around the idea of a series of authentication _providers_ that implement the interface **Authenticator**. This interface describes structures that have the ability to authenticate a user, save sessions, and verify credentials in a middleware. Currently, I only have basic user name and password authentication, but plan to fold in OAuth. 

## User Name and Password

This provider allows your application to keep user sessions based on user name and password credentials that you manaage in your own database infrastructure. To use this you will need a session storage, a structure to represent your user session, and routes for handling login. Here is an example. It features:

- Using the **sessions** package to setup a cookie store. You could also use PGStore for database sessions
- A custom session struct
- How to set the context key for where to put the session value when the middleware is done (for putting it in the context)
- Adding a custom handler in the event of an unexpected error (`WithErrorFunc`)
- How to specify where to send users when they are not logged in (`WithRedirectURL`)
- How to specify paths that are not verified in the authentication middleware

```go
type UserSession struct {
	Email string
}

gob.Register(&UserSession{})
sessionStore := sessions.NewCookieStore("mysessionkey")

auth := auth2.New(
	sessionStore,
	"my_app_session_name",
	"mysessionkey",
	auth2.WithContextKey("session"),
	auth2.WithErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		http.Redirect(w, r, "/error?message=" + err.Error(), http.StatusSeeOther)
	}),
	auth2.WithRedirectURL("/login"),
	auth2.WithExcludedPaths([]string{
		"/error",
		"/login",
	}),
)
```

Once you have your authentication configured, you are have an object that gives you the ability to add a middleware to routes, save sessions, and delete them. Here is a small example of using the middleware with the `mux` package.

```go
routes := []mux.Route{
	{Path: "GET /login", HandlerFunc: LoginPage},
	{Path: "POST /login", HandlerFunc: LoginAction},
}

routerConfig := mux.RouterConfig{
	Address:              config.Host,
	Debug:                Version == "development",
	ServeStaticContent:   true,
	StaticContentRootDir: "app",
	StaticContentPrefix:  "/static/",
	StaticFS:             appFS,
	Middlewares:          []mux.MiddlewareFunc{auth.Middleware},  // Authentication middleware
}

m := mux.SetupRouter(routerConfig, routes)
httpServer, quit := mux.SetupServer(routerConfig, m)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	// Do stuff to render the login page
}

func LoginAction(w http.ResponseWriter, r *http.Request) {
	// Validate the user, password...

	// Now create the session. Use the session struct we created above
	sessionValue := &UserSession{
		Email: user.Email,				// Pretend we got this somewhere above
	}

	err = auth.SaveSession(w, r, sessionValue)
	// handle error, or redirect if successfull
}
```

