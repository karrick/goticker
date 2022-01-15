package goticker

import (
	"context"
	"math/rand"
	"time"
)

// Sleep sleeps until either the context is canceled or the specified duration
// has elapsed. The difference between time.Sleep and this function is this
// function returns before the specified duration when the provided context is
// cancelled. It returns nil error when the sleep completed, or returns the
// context error if the context terminated before the duration elapsed.
func Sleep(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

// SleepMax calculates a random delay no longer than maxDuration, and sleeps
// until either the context is canceled or the delay has completed. It returns
// nil error when the sleep completed, or returns the context error if the
// context terminated before the duration elapsed.
func SleepMax(ctx context.Context, maxDuration time.Duration) error {
	d := time.Duration(rand.Float64() * float64(maxDuration.Nanoseconds()))
	return Sleep(ctx, d)
}

// SleepJitter will Sleep either until context.Context is cancelled or a
// pseudo-random length of time of base +/- jitter. It returns nil error when
// the sleep completed, or returns the context error if the context terminated
// before the duration elapsed. No checking is done whether base is less than
// jitter, whether adding base to a negative jitter results in a negative
// time.Duration, nor whether doubling jitter results in an overflow.
func SleepJitter(ctx context.Context, base, jitter time.Duration) error {
	nanos := jitter.Nanoseconds()

	// We want a pseudo-random number between [0, 1).
	//
	//     rand.Float64()
	//
	// Comments in the rand.Float64 implementation of The math/rand standard
	// library suggest using a clearer and simpler implementation, but the
	// library avoids the simpler implementation in order to maintain
	// compatibility with the random number sequence stream from Go 1. This
	// function does not have that requirement and we are free to take
	// advantage of that simplification and optimization:
	//
	//     rand.Float64() ::= float64(rand.Int63n(1<<53)) / (1<<53)
	//
	// Looking at implementation of rand.Int63n, and knowing that for our use
	// case n is always a power of 2, namely 1<<53, we can make the below
	// optimization, eliminating the chance of looping inside rand.Int63n.
	//
	//     rand.Int63n(n) ::= rand.Int63() & (n - 1)
	//
	i64 := rand.Int63() & ((1 << 53) - 1) // between [0, 1<<53), or technically [0, (1<<53)-1]
	f64 := float64(i64) / (1 << 53)       // between [0, 1.0)
	f64 *= 2.0                            // between [0, 2.0)
	f64 -= 1.0                            // between [-1.0, 1.0)
	f64 *= float64(nanos)                 // between [-nanos, nanos)
	d := base + time.Duration(f64)        // between [base-nanos, base+nanos)
	return Sleep(ctx, d)
}
