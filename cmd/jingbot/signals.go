package main

import (
	"context"
	"os"
	"os/signal"
)

// WithTerminationSignals handles OS signals and cancel() the given context.
func WithTerminationSignals(parent context.Context, sig ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(parent)

	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)

	// Call cancel func if signal encountered, or return if context ends.
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
			signal.Stop(c)
		}
	}()
	return ctx
}
