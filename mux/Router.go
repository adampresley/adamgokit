package mux

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/adampresley/adamgokit/auth"
	"github.com/rs/cors"
)

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

type RouterConfig struct {
	Address              string
	AuthConfig           *auth.AuthMiddlewareConfig
	Debug                bool
	HttpIdleTimeout      int
	HttpReadTimeout      int
	HttpWriteTimeout     int
	ServeStaticContent   bool
	StaticContentRootDir string
	StaticContentPrefix  string
	StaticFS             fs.FS
}

func SetupRouter(config RouterConfig, routes []Route) *http.ServeMux {
	var (
		staticFS      http.Handler
		excludedPaths []string
	)

	/*
	 * Ensure some sane defaults. Also panic on some things not being configured.
	 */
	if config.Address == "" {
		panic("router address cannot be blank.")
	}

	if config.StaticContentRootDir == "" {
		config.StaticContentRootDir = "app"
	}

	if config.StaticContentPrefix == "" {
		slog.Info("router StaticContentPrefix is blank. defaulting to '/static/'")
		config.StaticContentPrefix = "/static/"
	}

	if config.HttpIdleTimeout == 0 {
		config.HttpIdleTimeout = 120
	}

	if config.HttpReadTimeout == 0 {
		config.HttpReadTimeout = 120
	}

	if config.HttpWriteTimeout == 0 {
		config.HttpWriteTimeout = 120
	}

	m := http.NewServeMux()

	if config.ServeStaticContent {
		staticFS = http.FileServer(getStaticFileSystem(&config))
		m.Handle(fmt.Sprintf("GET %s", normalizeStaticContentPrefix(config.StaticContentPrefix)), staticFS)
	}

	if config.AuthConfig != nil {
		// Add provider endpoints to the exclusion list
		for _, authProvider := range config.AuthConfig.Providers {
			excludedPaths = append(excludedPaths, fmt.Sprintf("/%s", authProvider))
			excludedPaths = append(excludedPaths, fmt.Sprintf("/%s/callback", authProvider))
		}

		excludedPaths = append(excludedPaths, config.AuthConfig.ExcludedPaths...)
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
		if config.AuthConfig != nil {
			included := true

			for _, excluded := range excludedPaths {
				if strings.HasPrefix(route.Path, excluded) {
					included = false
					break
				}
			}

			if included {
				handler = config.AuthConfig.Middleware(handler)
			}

		}

		/*
		 * Wrap in any additional configured middlewares.
		 */
		for _, mw := range route.Middlewares {
			handler = mw(handler)
		}

		m.HandleFunc(route.Path, http.HandlerFunc(handler.ServeHTTP))
	}

	return m
}

func SetupServer(config RouterConfig, mux http.Handler) (*http.Server, chan os.Signal) {
	httpServer := &http.Server{
		Addr:         config.Address,
		WriteTimeout: time.Second * time.Duration(config.HttpWriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.HttpReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.HttpIdleTimeout),
		Handler:      cors.Default().Handler(mux),
	}

	go func() {
		slog.Info("starting HTTP server", slog.String("address", config.Address))
		err := httpServer.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			slog.Error("error starting HTTP server", slog.Any("error", err))
			os.Exit(-1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return httpServer, quit
}

func Shutdown(httpServer *http.Server) {
	httpContext, httpCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer httpCancel()

	_ = httpServer.Shutdown(httpContext)
	slog.Info("shut down complete")
}

func getStaticFileSystem(config *RouterConfig) http.FileSystem {
	if config.Debug {
		return http.FS(os.DirFS(config.StaticContentRootDir))
	}

	fsys, err := fs.Sub(config.StaticFS, config.StaticContentRootDir)

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
