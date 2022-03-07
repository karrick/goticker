package goticker

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// Config is a structure that controls behavior of a newly created
// Ticker.
type Config struct {
	// Callback specifies the callback function that will be invoked
	// by the Ticker. It will be called with the current time when it
	// is invoked.
	Callback func(time.Time)

	// Interval specifies how frequently the Callback should be
	// invoked.
	Interval time.Duration

	// Jitter controls how much time each interval should be
	// randomized over. Having a Jitter of 5 seconds means that for
	// each cycle a random value between -2.5 and +2.5 seconds will be
	// added to the time the go routine waits.
	Jitter time.Duration

	// Round allows rounding to the nearest specified unit of
	// time. For instance, if Interval were 24*time.Hour, and Round is
	// time.Minute, then the periodicity would be rounded to the
	// nearest minute on a 24 hour interval.
	Round time.Duration
}

// Ticker periodically invokes a callback function with the value of
// the current time. Allows callers to optionally specify whether
// invocations should occur at times that are rounded to the nearest
// duration interval. A Ticker will continue until its Stop method is
// invoked.
type Ticker struct {
	callback          func(time.Time)
	continueWhileZero uint32
}

// New spawns a goroutine that periodically invokes callback with the
// value of the current time. The periodicity is determined by
// interval. The requested duration must be greater than zero; if not,
// New will panic. The first invocation of the callback can be
// optionally rounded to the nearest duration interval by passing true
// for the round argument. Stop the ticker to release associated
// resources.
//
//     // Rotate logs every midnight...
//     logTicker, err := goticker.New(goticker.Config{
//         Interval: 24 * time.Hour,
//         Round:    time.Hour,
//         Callback: func(_ time.Time) {
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
//         Interval: time.Minute,
//         Jitter: 10*time.Second,
//         Callback: func(_ time.Time) {
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
		return nil, errors.New("cannot create Ticker when Callback omitted")
	}
	if c.Interval <= 0 {
		return nil, fmt.Errorf("cannot create Ticker when Interval is not greater than zero: %v", c.Interval)
	}
	if c.Jitter > 0 {
		if c.Interval < c.Jitter {
			return nil, fmt.Errorf("cannot create Ticker when Interval is smaller than Jitter: %v < %v", c.Interval, c.Jitter)
		}
		if c.Round > 0 {
			return nil, fmt.Errorf("cannot create Ticker when Jitter and Round are both non zero: %v and %v", c.Interval, c.Jitter)
		}
	}

	t := &Ticker{callback: c.Callback}

	if c.Jitter > 0 {
		go t.runJitter(c.Interval, c.Jitter)
	} else if c.Round > 0 {
		go t.runRound(c.Interval, c.Round)
	} else {
		go t.run(c.Interval)
	}

	return t, nil
}

// Stop will stop the Ticker preventing any further invocations of the
// Ticker's callback.
func (t *Ticker) Stop() {
	atomic.StoreUint32(&t.continueWhileZero, 1)
}

func (t *Ticker) run(interval time.Duration) {
	prev := time.Now()

	for atomic.LoadUint32(&t.continueWhileZero) == 0 {
		prev = prev.Add(interval)

		if d := time.Until(prev); d > 0 {
			time.Sleep(d)
			// POST: prev is current time
			t.callback(prev)
		} else {
			// Previous callback took longer than iteration; not sure
			// what time it is; therefore ask the system.
			t.callback(time.Now())
		}
	}
}

func (t *Ticker) runJitter(interval, jitter time.Duration) {
	prev := time.Now()

	for atomic.LoadUint32(&t.continueWhileZero) == 0 {
		prev = prev.Add(interval)
		next := prev.Add(plusOrMinus(jitter))

		if d := time.Until(next); d > 0 {
			time.Sleep(d)
			// POST: next is current time
			t.callback(next)
		} else {
			// Previous callback took longer than iteration; not sure
			// what time it is; therefore ask the system.
			t.callback(time.Now())
		}
	}
}

func (t *Ticker) runRound(interval, round time.Duration) {
	prev := time.Now()

	for atomic.LoadUint32(&t.continueWhileZero) == 0 {
		prev = prev.Add(interval)
		next := prev.Round(round)

		if d := time.Until(next); d > 0 {
			time.Sleep(d)
			// POST: next is current time
			t.callback(next)
		} else {
			// Previous callback took longer than iteration; not sure
			// what time it is; therefore ask the system.
			t.callback(time.Now())
		}
	}
}
