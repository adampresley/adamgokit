package getoptions

import (
	"context"
	"time"
)

type GetOptions struct {
	Context    context.Context
	Expiration time.Duration
	Timeout    time.Duration
}

type GetOption func(*GetOptions)

func WithContext(ctx context.Context) GetOption {
	return func(o *GetOptions) {
		o.Context = ctx
	}
}

func WithExpiration(expiration time.Duration) GetOption {
	return func(o *GetOptions) {
		o.Expiration = expiration
	}
}

func WithTimeout(timeout time.Duration) GetOption {
	return func(o *GetOptions) {
		o.Timeout = timeout
	}
}
