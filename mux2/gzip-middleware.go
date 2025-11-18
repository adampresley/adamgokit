package mux2

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipMiddlewareConfig struct {
	excludedPaths []string
}

type GzipMiddlewareConfigOption func(*GzipMiddlewareConfig)

func NewGzipMiddleware(options ...GzipMiddlewareConfigOption) func(next http.Handler) http.Handler {
	opts := &GzipMiddlewareConfig{
		excludedPaths: []string{},
	}

	for _, opt := range options {
		opt(opts)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			for _, excludedPath := range opts.excludedPaths {
				if strings.HasPrefix(r.URL.Path, excludedPath) {
					next.ServeHTTP(w, r)
					return
				}
			}

			gz := gzip.NewWriter(w)
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			gzw := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
			next.ServeHTTP(gzw, r)
		})
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithExcludedPaths(paths ...string) GzipMiddlewareConfigOption {
	return func(config *GzipMiddlewareConfig) {
		config.excludedPaths = append(config.excludedPaths, paths...)
	}
}
