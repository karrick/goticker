package goticker

import (
	"context"
	"math/rand"
	"time"
)

// Sleep sleeps until either the context is canceled or the delay has
// completed.
func Sleep(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

// SleepMax calculates a random delay no longer than maxDuration, and sleeps
// until either the context is canceled or the delay has completed.
func SleepMax(ctx context.Context, maxDuration time.Duration) error {
	d := time.Duration(rand.Float64() * float64(maxDuration.Nanoseconds()))
	return Sleep(ctx, d)
}
