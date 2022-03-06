package goticker

import (
	"math/rand"
	"time"
)

// plusOrMinus returns a pseudo-random number between [-jitter, jitter).
func plusOrMinus(jitter time.Duration) time.Duration {
	// We want a pseudo-random number between [-jitter, +jitter).
	//
	// To get that, we will need a pseudo-random number between [0, 1).
	//
	//     rand.Float64()
	//
	// Comments in the rand.Float64 implementation of The math/rand
	// standard library suggest using a clearer and simpler
	// implementation, but that library avoids the simpler
	// implementation in order to maintain compatibility with the
	// random number sequence stream from Go 1. This function does not
	// have that requirement and we are free to take advantage of that
	// simplification and optimization:
	//
	//     rand.Float64() ::= float64(rand.Int63n(1<<53)) / (1<<53)
	//
	// Looking at implementation of rand.Int63n, and knowing that for
	// our use case n is always a power of 2, namely 1<<53, we can
	// make the below optimization, eliminating the chance of looping
	// inside rand.Int63n.
	//
	//     rand.Int63n(n) ::= rand.Int63() & (n - 1)
	//
	i64 := rand.Int63() & ((1 << 53) - 1) // between [0, 1<<53), or technically [0, (1<<53)-1]
	f64 := float64(i64) / (1 << 53)       // between [0, 1.0)
	f64 *= 2.0                            // between [0, 2.0)
	f64 -= 1.0                            // between [-1.0, 1.0)
	f64 *= float64(jitter.Nanoseconds())  // between [-jitter, jitter)
	return time.Duration(f64)
}
