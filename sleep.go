package goticker

import (
	"context"
	"math/rand"
	"time"
)

// Sleep sleeps until either the context is canceled or the specified
// duration has elapsed. The difference between time.Sleep and this
// function is this function returns before the specified duration
// when the provided context is cancelled. It returns nil error when
// the sleep completed, or returns the context error if the context
// terminated before the duration elapsed.
func Sleep(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

// SleepMax calculates a random delay no longer than maxDuration, and
// sleeps until either the context is canceled or the delay has
// completed. It returns nil error when the sleep completed, or
// returns the context error if the context terminated before the
// duration elapsed.
func SleepMax(ctx context.Context, maxDuration time.Duration) error {
	d := time.Duration(rand.Float64() * float64(maxDuration.Nanoseconds()))
	return Sleep(ctx, d)
}

// SleepJitter will Sleep either until context.Context is cancelled or
// a pseudo-random length of time of base +/- jitter. It returns nil
// error when the sleep completed, or returns the context error if the
// context terminated before the duration elapsed. No checking is done
// whether base is less than jitter, whether adding base to a negative
// jitter results in a negative time.Duration, nor whether doubling
// jitter results in an overflow.
func SleepJitter(ctx context.Context, base, jitter time.Duration) error {
	return Sleep(ctx, base+plusOrMinus(jitter))
}
