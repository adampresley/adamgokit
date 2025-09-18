package calloptions

type CallOptions struct {
	Debug       bool
	Headers     map[string]string
	QueryParams map[string]string
}

type CallOption func(*CallOptions)

func WithCallHeaders(headers map[string]string) CallOption {
	return func(co *CallOptions) {
		co.Headers = headers
	}
}

func WithDebug(debug bool) CallOption {
	return func(co *CallOptions) {
		co.Debug = debug
	}
}

func WithQueryParams(params map[string]string) CallOption {
	return func(co *CallOptions) {
		co.QueryParams = params
	}
}
