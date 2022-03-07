package main

import (
	"fmt"
	"time"

	"github.com/karrick/goticker"
)

func main() {
	ticker1, err := goticker.New(goticker.Config{Interval: 2 * time.Second, Callback: func(t time.Time) {
		fmt.Println("TICKER1:", t, time.Since(t))
		time.Sleep(3)
	}})
	if err != nil {
		panic(err)
	}

	ticker2, err := goticker.New(goticker.Config{Interval: 10 * time.Second, Round: 10 * time.Second, Callback: func(t time.Time) {
		fmt.Println("TICKER2:", t, time.Since(t))
		time.Sleep(13)
	}})
	if err != nil {
		panic(err)
	}

	ticker3, err := goticker.New(goticker.Config{Interval: 10 * time.Second, Jitter: 3 * time.Second, Callback: func(t time.Time) {
		fmt.Println("TICKER3:", t, time.Since(t))
		time.Sleep(2)
	}})
	if err != nil {
		panic(err)
	}

	<-time.After(time.Minute)
	fmt.Printf("\n\ttest complete; stopping ticker...\n")

	ticker1.Stop()
	ticker2.Stop()
	ticker3.Stop()
}
