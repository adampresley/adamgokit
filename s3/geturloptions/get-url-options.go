package geturloptions

import (
	"context"
	"time"
)

type GetUrlOptions struct {
	Context    context.Context
	Expiration time.Duration
}

type GetUrlOption func(*GetUrlOptions)

func WithContext(ctx context.Context) GetUrlOption {
	return func(o *GetUrlOptions) {
		o.Context = ctx
	}
}

func WithExpiration(timeout time.Duration) GetUrlOption {
	return func(o *GetUrlOptions) {
		o.Expiration = timeout
	}
}
