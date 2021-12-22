package goticker

import (
	"context"
	"time"
)

// RunSleepLoop repeatedly invokes callback at a periodicity specified by
// duration, until the context is closed.
func RunSleepLoop(ctx context.Context, interval time.Duration, callback func(time.Time)) error {
	// One-time preperation before starting loop.
	ticker := time.NewTicker(interval)
	done := ctx.Done()
	now := time.Now()

	for {
		callback(now)

		// Wait until either context is done or ticker emits time over
		// channel.
		select {
		case <-done:
			ticker.Stop()
			return ctx.Err()
		case now = <-ticker.C:
			// Wake up.
		}
	}
}

// SleepRunLoop repeatedly invokes callback at a periodicity specified by
// duration, until the context is closed.
func SleepRunLoop(ctx context.Context, interval time.Duration, callback func(time.Time)) error {
	// One-time preperation before starting loop.
	ticker := time.NewTicker(interval)
	done := ctx.Done()
	var now time.Time

	for {
		// Wait until either context is done or ticker emits time over
		// channel.
		select {
		case <-done:
			ticker.Stop()
			return ctx.Err()
		case now = <-ticker.C:
			// Wake up.
		}

		callback(now)
	}
}
