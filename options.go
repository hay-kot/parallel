package parallel

import (
	"context"
	"runtime"
)

type options struct {
	maxProcs int
	ctx      context.Context
}

type optionsFunc func(*options)

func newOptions(opts ...optionsFunc) options {
	o := options{
		maxProcs: runtime.NumCPU(),
		ctx:      context.Background(),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

// WithMaxProcs sets the maximum number of parallel goroutines to use.
// within the ForEach function. If not set, the number of CPUs will be used.
func WithMaxProcs(maxProcs int) optionsFunc {
	return func(opts *options) {
		if maxProcs > 0 {
			opts.maxProcs = maxProcs
		}
	}
}

// WithContext sets the context to use for the parallel function. if not
// set, the context.Background() will be used.
func WithContext(ctx context.Context) optionsFunc {
	return func(opts *options) {
		opts.ctx = ctx
	}
}
