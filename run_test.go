package goticker

import (
	"context"
	"testing"
	"time"
)

func TestRunSleepLoop(t *testing.T) {
	t.Run("stops", func(t *testing.T) {
		const duration = time.Millisecond
		var prev time.Time

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			err := RunSleepLoop(ctx, duration, func(t time.Time) {
				prev = t
			})
			ensureError(t, err, context.Canceled.Error())
		}()

		// Allow to run for a few intervals, updating time stamp along the
		// way.
		<-time.After(5 * duration)

		cancel() // Canceling context should cause Run to terminate.
		stoppedAt := time.Now()

		// Wait a bit longer, and make sure prev has not updated (indicating
		// callback was invoked).
		<-time.After(100 * duration)

		if !stoppedAt.After(prev) {
			t.Errorf("stoppedAt: %v; prev: %v", stoppedAt, prev)
		}
	})
}

func TestSleepRunLoop(t *testing.T) {
	t.Run("stops", func(t *testing.T) {
		const duration = time.Millisecond

		ctx, cancel := context.WithCancel(context.Background())

		// Canceling context should prevent SleepRunLoop from invoking
		// callback.
		cancel()

		go func() {
			err := SleepRunLoop(ctx, duration, func(_ time.Time) {
				t.Error("should not have been invoked")
			})
			ensureError(t, err, context.Canceled.Error())
		}()
	})
}
