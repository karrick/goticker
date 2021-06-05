package goticker

import (
	"testing"
	"time"
)

func TestNewDurationLessThanOne(t *testing.T) {
	ensurePanic(t, "duration zero", "non-positive interval", func() {
		_ = New(0, false, func(_ time.Time) {})
	})
	ensurePanic(t, "duration negative", "non-positive interval", func() {
		_ = New(-time.Millisecond, false, func(_ time.Time) {})
	})
}

func TestStops(t *testing.T) {
	const duration = time.Millisecond

	var prev time.Time

	ticker := New(duration, false, func(t time.Time) {
		prev = t
	})

	// Let ticker run for a few intervals, updating prev along the way.
	<-time.After(5 * duration)

	ticker.Stop()
	stoppedAt := time.Now()

	// Let's wait a bit longer, and make sure prev has not updated (indicating
	// callback was invoked).
	<-time.After(100 * duration)
	if !stoppedAt.After(prev) {
		t.Errorf("stoppedAt: %v; prev: %v", stoppedAt, prev)
	}
}
