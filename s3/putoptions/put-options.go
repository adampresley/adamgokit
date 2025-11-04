package putoptions

import (
	"context"
	"time"
)

type PutOptions struct {
	ContentType string
	Context     context.Context
	Metadata    map[string]string
	Timeout     time.Duration
}

type PutOption func(*PutOptions)

func WithContentType(contentType string) PutOption {
	return func(o *PutOptions) {
		o.ContentType = contentType
	}
}

func WithContext(ctx context.Context) PutOption {
	return func(o *PutOptions) {
		o.Context = ctx
	}
}

func WithMetadata(metadata map[string]string) PutOption {
	return func(o *PutOptions) {
		o.Metadata = metadata
	}
}

func WithTimeout(timeout time.Duration) PutOption {
	return func(o *PutOptions) {
		o.Timeout = timeout
	}
}
