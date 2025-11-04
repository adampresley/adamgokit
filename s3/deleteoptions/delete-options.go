package deleteoptions

import (
	"context"
	"time"
)

type DeleteOptions struct {
	Context context.Context
	Timeout time.Duration
}

type DeleteOption func(*DeleteOptions)

func WithContext(ctx context.Context) DeleteOption {
	return func(opts *DeleteOptions) {
		opts.Context = ctx
	}
}

func WithTimeout(timeout time.Duration) DeleteOption {
	return func(opts *DeleteOptions) {
		opts.Timeout = timeout
	}
}
