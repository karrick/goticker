package main

import (
	"fmt"
	"time"

	"github.com/karrick/goticker"
)

func main() {
	cancel := goticker.NewTicker(time.Second, func(t time.Time) {
		fmt.Print(".")
		hr, min, sec := t.Clock()
		if sec == 0 /* && min == 0 */ /* && hr == 0 */ {
			fmt.Printf("\n\tthe time is %v: hr: %d; min: %d\n", t, hr, min)
		}
	})

	<-time.After(2 * time.Minute)
	fmt.Printf("\n\ttest complete; stopping ticker...\n")
	close(cancel)
}
