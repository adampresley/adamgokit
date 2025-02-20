# Mux

The **Mux** package provides structures and methods for removing some of the
boilerplate around setting up an HTTP server with routes, authentication
middlewares, and static routes. It makes use of the standard Go library
exclusively (requires Go 1.23.0 or higher).

Here is a short, basic example.

```go
handler := func(w http.ResponseWriter, r *http.Request) {
  httphelpers.TextOK(w, fmt.Sprintf("Hello %s", httphelpers.GetFromRequest[string](r, "name")))
}

routes := []mux.Route{
  {Path: "GET /", Handler: http.HandlerFunc(handler)},
}

routerConfig := mux.RouterConfig{
  Address: "localhost:8080",
  Debug:   true,
}

m := mux.SetupRouter(routerConfig, routes)
httpServer, quit := mux.SetupServer(routerConfig, m)

slog.Info("server started")

<-quit
mux.Shutdown(httpServer)
slog.Info("server stopped")
```

In this example, we first set up some routes that make use of HTTP handler
functions. Then we set up a configuration to describe how our HTTP server
is configured. Next we call the **SetupRouter** function to get a ServeMux
structure, and then pass that to **SetupServer**, which gives us
an HttpServer struct and a channel to wait for graceful shutdown.

## Routes

A **route** is simply a structure that defines the handler and any middlewares
for a given verb and path. Here is what that looks like:

```go
mux.Route{
  Path: "GET /about",
  Handler: http.HandlerFunc(SomeHandlerFunc),
  Middlewares: []func(http.Handler) http.Handler{
    SomeMiddlewareFunc,
  }
}
```

_Handler_ and _Path_ are the only fields required. _Middlewares_ is optional.

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

- Per route
- Built-in authentication middlewraes via configuration

### Per-route

Here is an example of adding a middleware to a single route.

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

routerConfig := mux.RouterConfig{
  Address:    "localhost:8080",
  AuthConfig: &auth.AuthMiddlewareConfig{
    ExcludedPaths: []string{"/login", "/logout"},
    Middleware: authMiddleware,
  },
  Debug:      true,
}
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

// More code and stuff
routerConfig := mux.RouterConfig{
  Address:             "localhost:8080",
  Debug:               true,
  ServeStaticContent:  true,
  StaticContextPrefix: "/static/",
  StaticFS:            appFS,
}
```
