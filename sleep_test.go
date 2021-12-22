package goticker

import (
	"context"
	"testing"
	"time"
)

// Epsilon is a Time difference less than which is ignored.
const epsilon = 100 * time.Microsecond

func ensureDurationApproximately(tb testing.TB, got, want time.Duration) {
	tb.Helper()

	// | got - want | < epsilon
	diff := got - want

	if diff > 0 {
		if diff > epsilon {
			tb.Errorf("difference too large: %v > %v", got, want)
		}
	} else {
		if -diff > epsilon {
			tb.Errorf("difference too large: %v > %v", got, want)
		}
	}
}

func measure(callback func()) time.Duration {
	start := time.Now()
	callback()
	return time.Since(start)
}

func TestSleep(t *testing.T) {
	t.Run("does not wait if context already closed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel context before invoking Sleep

		diff := measure(func() { Sleep(ctx, time.Second) })

		ensureDurationApproximately(t, diff, 0)
	})

	t.Run("waits required amount of time", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		diff := measure(func() { Sleep(ctx, time.Millisecond) })

		ensureDurationApproximately(t, diff, time.Millisecond)
	})
}

func TestSleepMax(t *testing.T) {
	t.Run("does not wait if context already closed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel context before invoking Sleep

		diff := measure(func() { SleepMax(ctx, time.Second) })

		ensureDurationApproximately(t, diff, 0)
	})

	t.Run("waits no longer than specified amount of time", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		diff := measure(func() { Sleep(ctx, time.Millisecond) })

		if got, want := diff, time.Millisecond+epsilon; diff > want {
			t.Errorf("difference too large: %v > %v", got, want)
		}
	})
}
