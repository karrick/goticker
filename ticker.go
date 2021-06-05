package goticker

import (
	"sync/atomic"
	"time"
)

// Ticker periodically invokes a callback function with the value of the current
// time. Allows callers to optionally specify whether invocations should occur
// at times that are rounded to the nearest duration interval. A Ticker will
// continue until its Stop method is invoked.
type Ticker struct {
	callback func(time.Time)
	duration time.Duration
}

// New spawns a go routine that periodically invokes callback every duration
// nanoseconds, optionally rounded to the nearest duration interval.
//
//     func main() {
//         ticker1 := goticker.New(5*time.Second, false, func(t time.Time) {
//             fmt.Println(t, false)
//             time.Sleep(1)
//         })
//         ticker2 := goticker.New(5*time.Second, true, func(t time.Time) {
//             fmt.Println(t, true)
//             time.Sleep(1)
//         })
//
//         <-time.After(time.Minute)
//         fmt.Printf("\n\ttest complete; stopping ticker...\n")
//
//         ticker1.Stop()
//         ticker2.Stop()
//     }
func New(duration time.Duration, round bool, callback func(time.Time)) *Ticker {
	t := &Ticker{duration: duration, callback: callback}

	// By sending duration to methods on stack, we elide an atomic load.
	if round {
		go t.runRound(duration)
	} else {
		go t.run(duration)
	}

	return t
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
