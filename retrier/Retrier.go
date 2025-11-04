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
		Delay:                time.Second * 5,
		MaxAttempts:          3,
		MaxJitter:            300,
		MaxJitterMeasurement: time.Millisecond,
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

		sleepDuration := opts.Delay + (time.Duration(rand.Intn(opts.MaxJitter)) * opts.MaxJitterMeasurement)
		time.Sleep(sleepDuration)
	}

	return fmt.Errorf("failed to execute function after %d attempts: %w", opts.MaxAttempts, err)
}

type Options struct {
	Delay                time.Duration
	MaxAttempts          int
	MaxJitter            int
	MaxJitterMeasurement time.Duration
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

func WithMaxJitter(jitter int) Option {
	return func(o *Options) {
		o.MaxJitter = jitter
	}
}

func WithMaxJitterMeasurement(measurement time.Duration) Option {
	return func(o *Options) {
		o.MaxJitterMeasurement = measurement
	}
}
