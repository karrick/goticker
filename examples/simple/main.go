package main

import (
	"fmt"
	"time"

	"github.com/karrick/goticker"
)

func main() {
	cancel1 := goticker.NewTicker(5*time.Second, false, func(t time.Time) {
		fmt.Println(t, false)
		time.Sleep(1)
	})

	cancel2 := goticker.NewTicker(5*time.Second, true, func(t time.Time) {
		fmt.Println(t, true)
		time.Sleep(1)
	})

	<-time.After(time.Minute)
	fmt.Printf("\n\ttest complete; stopping ticker...\n")

	close(cancel1)
	close(cancel2)
}
