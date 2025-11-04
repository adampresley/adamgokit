package listoptions

import (
	"context"
	"time"

	"github.com/adampresley/adamgokit/s3/geturloptions"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type FilterFunc func(obj types.Object) bool

type ListOptions struct {
	Context           context.Context
	ContinuationToken string
	Filter            FilterFunc
	GetAll            bool
	GetUrls           bool
	GetUrlOptions     *geturloptions.GetUrlOptions
	Timeout           time.Duration
}

type ListOption func(*ListOptions)

func WithContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}

func WithContinuationToken(token string) ListOption {
	return func(o *ListOptions) {
		o.ContinuationToken = token
	}
}

func WithFilter(filter FilterFunc) ListOption {
	return func(o *ListOptions) {
		o.Filter = filter
	}
}

func WithGetAll() ListOption {
	return func(o *ListOptions) {
		o.GetAll = true
	}
}

func WithGetUrls() ListOption {
	return func(o *ListOptions) {
		o.GetUrls = true
	}
}

func WithGetUrlOptions(options ...geturloptions.GetUrlOption) ListOption {
	return func(o *ListOptions) {
		opts := &geturloptions.GetUrlOptions{
			Context:    context.Background(),
			Expiration: time.Second * 5,
		}

		for _, opt := range options {
			opt(opts)
		}

		o.GetUrlOptions = opts
	}
}

func WithTimeout(d time.Duration) ListOption {
	return func(o *ListOptions) {
		o.Timeout = d
	}
}
