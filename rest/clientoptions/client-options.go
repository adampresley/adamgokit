package clientoptions

import (
	"net/http"

	"github.com/adampresley/adamgokit/httphelpers"
)

type ClientOptions struct {
	BaseURL    string
	Debug      bool
	Headers    map[string]string
	HttpClient httphelpers.HttpClient
}

type ClientOption func(*ClientOptions)

func New(baseURL string, options ...ClientOption) *ClientOptions {
	result := &ClientOptions{
		BaseURL:    baseURL,
		HttpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(result)
	}

	return result
}

func WithDebug(debug bool) ClientOption {
	return func(s *ClientOptions) {
		s.Debug = debug
	}
}

func WithHeaders(headers map[string]string) ClientOption {
	return func(s *ClientOptions) {
		s.Headers = headers
	}
}

func WithHttpClient(client httphelpers.HttpClient) ClientOption {
	return func(s *ClientOptions) {
		s.HttpClient = client
	}
}
