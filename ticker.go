package goticker

import "time"

// NewTicker spawns a go routine that periodically invokes callback every
// duration nanoseconds. It returns a channel that when closed, stops the go
// routine.
//
//     func main() {
//         cancel := goticker.NewTicker(time.Second, func(t time.Time) {
//             fmt.Print(".")
//             hr, min, sec := t.Clock()
//             if sec == 0 /* && min == 0 */ /* && hr == 0 */ {
//                 fmt.Printf("\n\tthe time is %v: hr: %d; min: %d\n", t, hr, min)
//             }
//         })
//
//         <-time.After(2 * time.Minute)
//         fmt.Printf("\n\ttest complete; stopping ticker...\n")
//         close(cancel)
//     }
func NewTicker(duration time.Duration, callback func(time.Time)) chan struct{} {
	// Create buffered channel because while caller only needs to close it to
	// stop the ticker, they might send to it instead, and we don't want that
	// send action to block their go routine while this go routine is waiting
	// for callback to return.
	cancel := make(chan struct{}, 1)

	go func() {
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-cancel:
				ticker.Stop()
				return
			case t := <-ticker.C:
				callback(t)
			}
		}
	}()

	return cancel
}
