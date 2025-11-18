package mux2

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adampresley/adamgokit/waiter"
	"github.com/rs/cors"
)

/*
A Route defines a single handler for a single endpoint. You have a choice
of using the http.Handler interface or the direct http.HandlerFunc. You
may also provide an optional slice of middlewares that will be automatically
wrapped around your handler.

For example, using http.HandlerFunc might look like:

	func homePage(w http.ResponseWriter, r *http.Request) {
	  fmt.Fprintf(w, "This is a test")
	}

	routes := []mux.Route{
	  {Path: "GET /", HandlerFunc: homePage},
	}
*/
type Route struct {
	Path        string
	Handler     http.Handler
	HandlerFunc http.HandlerFunc
	Middlewares []MiddlewareFunc
}

type Router struct {
	Mux    *http.ServeMux
	Server *http.Server

	opts        *routerConfig
	shutdownCtx context.Context
	stopApp     context.CancelFunc
}

/*
Setup an HTTP mux and server for handling HTTP requests.
*/
func Setup(config MuxConfig, routes []Route, shutdownCtx context.Context, stopApp context.CancelFunc, options ...RouterOption) *Router {
	opts := &routerConfig{
		address:              config.GetHost(),
		authConfig:           nil,
		cors:                 cors.AllowAll(),
		debug:                false,
		faviconPath:          "/static/favicons",
		httpIdleTimeout:      2 * time.Minute,
		httpReadTimeout:      1 * time.Minute,
		httpWriteTimeout:     1 * time.Minute,
		letsEncryptConfig:    nil,
		middlewares:          []MiddlewareFunc{},
		serveStaticContent:   false,
		staticContentRootDir: "",
		staticContentPrefix:  "",
		staticFS:             nil,
		useGzipForStaticFS:   false,
	}

	for _, opt := range options {
		opt(opts)
	}

	validateConfig(opts)

	m := setupMux(routes, opts)
	s := setupServer(opts, m)

	result := &Router{
		Mux:    m,
		Server: s,

		opts:        opts,
		shutdownCtx: shutdownCtx,
		stopApp:     stopApp,
	}

	return result
}

/*
Starts the HTTP server. This method blocks until the server is stopped.
*/
func (r *Router) Start() {
	var (
		err error
	)

	slog.Info("starting HTTP server", slog.String("address", r.opts.address))

	if r.opts.letsEncryptConfig != nil {
		err = r.Server.ListenAndServeTLS("", "")
	} else {
		err = r.Server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		slog.Error("error starting HTTP server", slog.Any("error", err))
		os.Exit(-1)
	}

	go func() {
		<-waiter.Wait()
		r.stopApp()
	}()

	<-r.shutdownCtx.Done()
}

func validateConfig(config *routerConfig) {
	if config.address == "" {
		panic("router address cannot be blank.")
	}

	if config.serveStaticContent {
		if config.staticContentRootDir == "" {
			slog.Info("static content root directory is blank. defaulting to 'app'")
			config.staticContentRootDir = "app"
		}

		if config.staticContentPrefix == "" {
			slog.Info("static content prefix is blank. defaulting to '/static/'")
			config.staticContentPrefix = "/static/"
		}
	}
}

func setupMux(routes []Route, opts *routerConfig) *http.ServeMux {
	var (
		staticFS      http.Handler
		excludedPaths []string
	)

	m := http.NewServeMux()

	if opts.serveStaticContent {
		staticFS = http.FileServer(getStaticFileSystem(opts))
		var wrappedStaticFS http.Handler = staticFS

		if opts.useGzipForStaticFS {
			wrappedStaticFS = NewGzipMiddleware()(wrappedStaticFS)
		}

		m.Handle(fmt.Sprintf("GET %s", normalizeStaticContentPrefix(opts.staticContentPrefix)), wrappedStaticFS)
	}

	if opts.authConfig != nil {
		// Add provider endpoints to the exclusion list
		for _, authProvider := range opts.authConfig.Providers {
			excludedPaths = append(excludedPaths, fmt.Sprintf("/%s", authProvider))
			excludedPaths = append(excludedPaths, fmt.Sprintf("/%s/callback", authProvider))
		}

		excludedPaths = append(excludedPaths, opts.authConfig.ExcludedPaths...)
	}

	for _, route := range routes {
		var handler http.Handler

		if route.HandlerFunc != nil {
			handler = http.HandlerFunc(route.HandlerFunc)
		}

		if route.Handler != nil {
			handler = route.Handler
		}

		/*
		 * If we have an auth configuration, and the path isn't excluded,
		 * wrap in the auth middleware.
		 */
		if opts.authConfig != nil {
			included := true

			for _, excluded := range excludedPaths {
				if strings.HasPrefix(route.Path, excluded) {
					included = false
					break
				}
			}

			if included {
				handler = opts.authConfig.Middleware(handler)
			}

		}

		/*
		 * Wrap in any additional router-configured middlewares.
		 */
		for _, mw := range opts.middlewares {
			handler = mw(handler)
		}

		/*
		 * Wrap in any additional route-configured middlewares.
		 */
		for _, mw := range route.Middlewares {
			handler = mw(handler)
		}

		/* Wrap in gzip middleware if enabled */
		if opts.useGzip {
			handler = NewGzipMiddleware(WithExcludedPaths(opts.gzipExcludedPaths...))(handler)
		}

		m.HandleFunc(route.Path, http.HandlerFunc(handler.ServeHTTP))
	}

	return m
}

func setupServer(opts *routerConfig, m *http.ServeMux) *http.Server {
	var (
		tlsConfig *tls.Config
		server    *http.Server
	)

	if opts.letsEncryptConfig != nil {
		tlsConfig = startCertManager(*opts.letsEncryptConfig)
	}

	server = &http.Server{
		Addr:         opts.address,
		WriteTimeout: opts.httpWriteTimeout,
		ReadTimeout:  opts.httpReadTimeout,
		IdleTimeout:  opts.httpIdleTimeout,
		Handler:      opts.cors.Handler(m),
		TLSConfig:    tlsConfig,
	}

	return server
}

func getStaticFileSystem(opts *routerConfig) http.FileSystem {
	if opts.debug {
		return http.FS(os.DirFS(opts.staticContentRootDir))
	}

	fsys, err := fs.Sub(opts.staticFS, opts.staticContentRootDir)

	if err != nil {
		slog.Error("error loading static asset filesystem", slog.Any("error", err))
		os.Exit(-1)
	}

	return http.FS(fsys)
}

func normalizeStaticContentPrefix(prefix string) string {
	result := ""

	if !strings.HasPrefix(prefix, "/") {
		result += "/"
	}

	result += prefix

	if !strings.HasSuffix(result, "/") {
		result += "/"
	}

	return result
}
