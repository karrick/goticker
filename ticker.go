package goticker

import "time"

// NewTicker spawns a go routine that periodically invokes callback every
// duration nanoseconds. It returns a channel that when closed, stops the go
// routine.
//
//      func main() {
//          cancel1 := goticker.NewTicker(5*time.Second, false, func(t time.Time) {
//              fmt.Println(t, false)
//              time.Sleep(1)
//          })
//
//          cancel2 := goticker.NewTicker(5*time.Second, true, func(t time.Time) {
//              fmt.Println(t, true)
//              time.Sleep(1)
//          })
//
//          <-time.After(time.Minute)
//          fmt.Printf("\n\ttest complete; stopping ticker...\n")
//
//          close(cancel1)
//          close(cancel2)
//      }
func NewTicker(duration time.Duration, round bool, callback func(time.Time)) chan struct{} {
	// Create buffered channel because while caller only needs to close it to
	// stop the ticker, they might send to it instead, and we don't want that
	// send action to block their go routine while this go routine is waiting
	// for callback to return.
	cancel := make(chan struct{}, 1)

	go func() {
		prev := time.Now()

		for {
			// Next time to wake up should be duration nanoseconds after
			// previous wake up time, ignoring how long previous callback took.
			next := prev.Add(duration)
			if round {
				next = next.Round(duration)
			}
			time.Sleep(next.Sub(prev))

			// Non-blocking receive from cancel channel.
			select {
			case _ = <-cancel:
				// Channel has been closed or received from.
				return
			default:
				// Default case when channel remains open yet have received
				// nothing.
			}

			prev = time.Now()
			callback(prev)
		}
	}()

	return cancel
}
