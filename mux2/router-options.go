package mux2

import (
	"context"
	"io/fs"
	"time"

	"github.com/adampresley/adamgokit/auth"
	"github.com/rs/cors"
)

type routerConfig struct {
	address              string
	authConfig           *auth.AuthMiddlewareConfig
	cors                 *cors.Cors
	debug                bool
	faviconPath          string
	httpIdleTimeout      time.Duration
	httpReadTimeout      time.Duration
	httpWriteTimeout     time.Duration
	letsEncryptConfig    *LetsEncryptConfig
	middlewares          []MiddlewareFunc
	serveStaticContent   bool
	shutdownCtx          context.Context
	staticContentRootDir string
	staticContentPrefix  string
	staticFS             fs.FS
	useGzip              bool
	gzipExcludedPaths    []string
	useGzipForStaticFS   bool
}

type RouterOption func(r *routerConfig)

func WithAuth(config *auth.AuthMiddlewareConfig) RouterOption {
	return func(r *routerConfig) {
		r.authConfig = config
	}
}

func WithCors(options cors.Options) RouterOption {
	return func(r *routerConfig) {
		r.cors = cors.New(options)
	}
}

func WithDebug(enable bool) RouterOption {
	return func(r *routerConfig) {
		r.debug = enable
	}
}

func FaviconPath(path string) RouterOption {
	return func(r *routerConfig) {
		r.faviconPath = path
	}
}

func WithIdleTimeout(timeout time.Duration) RouterOption {
	return func(r *routerConfig) {
		r.httpIdleTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) RouterOption {
	return func(r *routerConfig) {
		r.httpReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) RouterOption {
	return func(r *routerConfig) {
		r.httpWriteTimeout = timeout
	}
}

func WithLetsEncrypt(config *LetsEncryptConfig) RouterOption {
	return func(r *routerConfig) {
		r.letsEncryptConfig = config
	}
}

func WithMiddlewares(middlewares ...MiddlewareFunc) RouterOption {
	return func(r *routerConfig) {
		r.middlewares = append(r.middlewares, middlewares...)
	}
}

func WithStaticContent(rootDir, prefix string, fs fs.FS) RouterOption {
	return func(r *routerConfig) {
		r.serveStaticContent = true
		r.staticContentRootDir = rootDir
		r.staticContentPrefix = prefix
		r.staticFS = fs
	}
}

func UseGzip(excludedPaths ...string) RouterOption {
	return func(r *routerConfig) {
		r.useGzip = true
		r.gzipExcludedPaths = excludedPaths
	}
}

func UseGzipForStaticFiles() RouterOption {
	return func(r *routerConfig) {
		r.useGzipForStaticFS = true
	}
}
