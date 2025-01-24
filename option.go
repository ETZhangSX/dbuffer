package dbuffer

import (
	"time"
)

type Options struct {
	interval time.Duration
}

func (o *Options) init() {
	if o.interval == 0 {
		o.interval = DefaultInterval
	}
}

const (
	Off = -1
	// DefaultInterval is used when Option.interval is not set.
	DefaultInterval = time.Second * 5
)

type Option func(*Options)

// WithInterval set the interval for double buffer
func WithInterval(d time.Duration) Option {
	return func(opts *Options) {
		if d < 0 {
			opts.interval = -1
			return
		}
		opts.interval = d
	}
}
