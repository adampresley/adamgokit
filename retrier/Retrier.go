package retrier

import (
	"fmt"
	"math/rand"
	"time"
)

func Retry(fn func() error, options ...Option) error {
	var (
		err error
	)

	opts := &Options{
		Delay:       time.Second * 2,
		MaxAttempts: 3,
		MaxJitter:   300 * time.Millisecond,
	}

	for _, opt := range options {
		opt(opts)
	}

	attempt := 0

	for attempt < opts.MaxAttempts {
		err = fn()

		if err == nil {
			return nil
		}

		sleepDuration := (opts.Delay * time.Duration(attempt+1)) * time.Duration(rand.Intn(int(opts.MaxJitter)))
		time.Sleep(sleepDuration)
	}

	return fmt.Errorf("failed to execute function after %d attempts: %w", opts.MaxAttempts, err)
}

type Options struct {
	Delay       time.Duration
	MaxAttempts int
	MaxJitter   time.Duration
}

type Option func(o *Options)

func WithDelay(delay time.Duration) Option {
	return func(o *Options) {
		o.Delay = delay
	}
}

func WithMaxAttempts(attempts int) Option {
	return func(o *Options) {
		o.MaxAttempts = attempts
	}
}

func WithMaxJitter(jitter time.Duration) Option {
	return func(o *Options) {
		o.MaxJitter = jitter
	}
}
