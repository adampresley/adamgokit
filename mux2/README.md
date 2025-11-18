# Mux

The **Mux** package provides structures and methods for removing some of the
boilerplate around setting up an HTTP server with routes, authentication
middlewares, and static routes. It makes use of the standard Go library
exclusively (requires Go 1.25.0 or higher).

Here is a short, basic example.

```go
type Config struct {
  mux2.MuxConfig
}

// Load your config using something like Configinator

shutdownCtx, stopApp := context.WithCancel(context.Background())

handler := func(w http.ResponseWriter, r *http.Request) {
  httphelpers.TextOK(w, fmt.Sprintf("Hello %s", httphelpers.GetFromRequest[string](r, "name")))
}

routes := []mux.Route{
  {Path: "GET /", HandlerFunc: handler},
}

mux := mux2.Setup(
  &config,
  routes,
  shutdownCtx,
  stopApp,
)

slog.Info("server started")
mux.Start()
slog.Info("server stopped")
```

In this example, we need and setup the following:

- A cancellable context
- A config struct (that implements **mux2.MuxConfig**
- Some routes with handlers
- Setup our mux
- Start the HTTP server (this blocks until the context is cancelled. terminate signals are built in)

## Routes

A **route** is simply a structure that defines the handler and any middlewares
for a given verb and path. Here is what that looks like:

```go
mux2.Route{
  Path: "GET /about",
  HandlerFunc: http.HandlerFunc(SomeHandlerFunc),
  Middlewares: []mux.MiddlewareFunc{
    SomeMiddlewareFunc,
  }
}
```

_HandlerFunc_ (or _Handler_) and _Path_ are the only fields required. _Middlewares_ is optional.

## Middlewares

Middleware functions allow you to run a method prior to a handler servicing
the request. These functions can alter the request context, validate
information, and more. Their signature looks like this.

```go
func middlewareFunc(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Do stuff
    next.ServeHTTP(w, r)
  })
}
```

This library supports two ways to intercept requests with middlewares:

- Router-level
- Per route
- Built-in authentication middlewares via configuration

### Router-level

Router-level middlewares are applied to all routes. This is useful in situations 
such as logging. Here is a simple example:

```go
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request path: %s\n", r.Path)
	})
}

mux := mux2.Setup(
  &config,
  routes,
  shutdownCtx,
  stopApp,

  mux2.WithMiddlewares(loggingMiddleware),
)
```

### Per-route

Per-route middlewares are middlewares that are only applied to a single route.
This is useful if you wish to intercept a single request to perform some action.

```go
func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do something here.
	})
}

routes := []mux.Route{
	{Path: "GET /about", HandlerFunc: aboutHandler, Middlewares: []mux.MiddlewareFunc{testMiddleware}},
}
```

### Authentication middleware

The router will automatically register an authentication middleware of your
choice if it is set up in the router config. It provides a way to configure
paths that should be ignored by the middleware if they match as a prefix.

For example:

```go
func authMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Do some auth stuff
    next.ServeHTTP(w, r)
  })
}

authConfig := &auth.AuthMiddlewareConfig{
  ExcludedPaths: []string{"/login", "/logout"},
  Middleware: authMiddleware,
}

mux := mux2.Setup(
  &config,
  routes,
  shutdownCtx,
  stopApp,

  mux2.WithAuth(authConfig),
)
```

In the above example we tell our router that if the path doesn't start with _/login_
or _/logout_, run requests through the **authMiddleware**.

## Static Assets

To serve static assets you will need to do four things:

1. Create a folder where static assets live. For example: `/app`
2. Create a sub-folder that matches what your static asset path should be:
   `/app/static`
3. Add an embed to your application
4. Add static content configuration to the router

All of your static assets should live inside the sub-folder `static`. Here is
a sample of the code for setting up static assets.

```go
import (
  "embed"
)

var (
  //go:embed app
  appFS embed.FS
)

mux := mux2.Setup(
  &config,
  routes,
  shutdownCtx,
  stopApp,

  mux2.WithStaticContent("app", "/static/", appFS),
  mux2.WithGzipForStaticFiles(),
)
```

