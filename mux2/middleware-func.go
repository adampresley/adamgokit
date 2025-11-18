package mux2

import "net/http"

/*
Defines a type for a middleware function. It must look like: `func(http.Handler) http.Handler {}`.
Here is an example:

	func logMiddleware(next http.Handler) http.Handler {
	   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	      slog.Info("running path", "path", r.URL.Path)
	      next.ServeHTTP(w, r)
	   })
	}
*/
type MiddlewareFunc func(http.Handler) http.Handler
