package recovery

import (
	"context"
)

var (
	defaultOptions = &options{
		recoveryHandlerFunc: nil,
	}
)

type options struct {
	recoveryHandlerFunc HandlerFuncContext
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

// WithRecoveryHandler customizes the function for recovering from a panic.
func WithRecoveryHandler(f HandlerFunc) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = HandlerFuncContext(func(ctx context.Context, p interface{}) error {
			return f(p)
		})
	}
}

// WithRecoveryHandlerContext customizes the function for recovering from a panic.
func WithRecoveryHandlerContext(f HandlerFuncContext) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = f
	}
}
