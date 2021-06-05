package goticker

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// Config is a structure that controls behavior of a newly created Ticker.
type Config struct {
	// Callback specifies the callback function that will be invoked by the
	// Ticker. It will be called with the current time when it is invoked.
	Callback func(time.Time)

	// Duration specifies how frequently the Callback should be invoked.
	Duration time.Duration

	// Round controls whether the ticks should occur at time intervals that are
	// rounded on the tick duration. For example, assume Duration is time.Minute
	// and the Ticker is created at 13 seconds after the current minute. When
	// Round is false, then Callback will be invoked every minute, 13 seconds
	// after the minute started. When Round is true, however, then Callback will
	// be invoked every minute, at 0 seconds after the minute started, and the
	// first Callback invocation would happen 47 seconds after New was invoked.
	Round bool
}

// Ticker periodically invokes a callback function with the value of the current
// time. Allows callers to optionally specify whether invocations should occur
// at times that are rounded to the nearest duration interval. A Ticker will
// continue until its Stop method is invoked.
type Ticker struct {
	callback func(time.Time)
	duration time.Duration
}

// New spawns a go routine that periodically invokes callback with the value of
// the current time. The periodicity is determined by duration. The requested
// duration must be greater than zero; if not, New will panic. The first
// invocation of the callback can be optionally rounded to the nearest duration
// interval by passing true for the round argument. Stop the ticker to release
// associated resources.
//
//     // Rotate logs every midnight...
//     logTicker, err := goticker.New(goticker.Config{
//         Duration: 24 * time.Hour,
//         Round:    true,
//         Callback: func(t time.Time) {
//             logger.Rotate()
//         }})
//     if err != nil {
//         panic(err) // TODO: handle appropriately
//     }
//
//     // some time later...
//     logTicker.Stop()
//
//     // Emit metrics every minute...
//     metricTicker, err := goticker.New(goticker.Config{
//         Duration: time.Minute,
//         Callback: func(t time.Time) {
//             metrics.Emit()
//         }})
//     if err != nil {
//         panic(err) // TODO: handle appropriately
//     }
//
//     // some time later...
//     metricTicker.Stop()
func New(c Config) (*Ticker, error) {
	if c.Callback == nil {
		return nil, errors.New("callback omitted")
	}
	if c.Duration <= 0 {
		return nil, fmt.Errorf("non-positive interval for Ticker: %v", c.Duration)
	}

	t := &Ticker{duration: c.Duration, callback: c.Callback}

	// By sending duration to methods on stack, we elide an atomic load.
	if c.Round {
		go t.runRound(c.Duration)
	} else {
		go t.run(c.Duration)
	}

	return t, nil
}

// Stop will stop the Ticker preventing any further invocations of the Ticker's
// callback.
func (t *Ticker) Stop() {
	atomic.StoreInt64((*int64)(&t.duration), 0)
}

func (t *Ticker) run(duration time.Duration) {
	prev := time.Now()

	for {
		// Next time to wake up should be duration nanoseconds after
		// previous wake up time, ignoring how long previous callback took.
		time.Sleep(prev.Add(duration).Sub(prev))

		if duration = time.Duration(atomic.LoadInt64((*int64)(&t.duration))); duration == 0 {
			return
		}

		prev = time.Now()
		t.callback(prev)
	}
}

func (t *Ticker) runRound(duration time.Duration) {
	prev := time.Now()

	for {
		// Next time to wake up should be duration nanoseconds after
		// previous wake up time, ignoring how long previous callback took.
		time.Sleep(prev.Add(duration).Round(duration).Sub(prev))

		if duration = time.Duration(atomic.LoadInt64((*int64)(&t.duration))); duration == 0 {
			return
		}

		prev = time.Now()
		t.callback(prev)
	}
}
