package main

import (
	"fmt"
	"time"

	"github.com/karrick/goticker"
)

func main() {
	ticker1 := goticker.New(5*time.Second, false, func(t time.Time) {
		fmt.Println(t, false)
		time.Sleep(1)
	})
	ticker2 := goticker.New(5*time.Second, true, func(t time.Time) {
		fmt.Println(t, true)
		time.Sleep(1)
	})

	<-time.After(time.Minute)
	fmt.Printf("\n\ttest complete; stopping ticker...\n")

	ticker1.Stop()
	ticker2.Stop()
}
