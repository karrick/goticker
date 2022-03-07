package goticker

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("interval", func(t *testing.T) {
		t.Run("0", func(t *testing.T) {
			_, err := New(Config{Callback: func(time.Time) {}})
			ensureError(t, err, "not greater than zero")
		})
		t.Run("negative", func(t *testing.T) {
			_, err := New(Config{Interval: -time.Millisecond, Callback: func(time.Time) {}})
			ensureError(t, err, "not greater than zero")
		})
	})
	t.Run("callback omitted", func(t *testing.T) {
		_, err := New(Config{Interval: time.Millisecond})
		ensureError(t, err, "Callback omitted")
	})
}

func TestStops(t *testing.T) {
	const duration = time.Millisecond

	var prev time.Time

	ticker, err := New(Config{
		Interval: duration,
		Callback: func(t time.Time) {
			prev = t
		},
	})
	ensureError(t, err)

	// Let ticker run for a few intervals, updating prev along the
	// way.
	<-time.After(5 * duration)
	ticker.Stop()

	// Let's wait a bit longer, and make sure prev has not updated
	// (indicating callback was invoked).
	<-time.After(10 * duration)

	stoppedAt := time.Now()
	if !stoppedAt.After(prev) {
		t.Errorf("stoppedAt: %v; prev: %v", stoppedAt, prev)
	}
}
