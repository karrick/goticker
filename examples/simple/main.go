package main

import (
	"fmt"
	"time"

	"github.com/karrick/goticker"
)

func main() {
	ticker1, err := goticker.New(goticker.Config{Duration: 5 * time.Second, Callback: func(t time.Time) {
		fmt.Println(t, false)
		time.Sleep(1)
	}})
	if err != nil {
		panic(err)
	}

	ticker2, err := goticker.New(goticker.Config{
		Duration: 5 * time.Second,
		Round:    true,
		Callback: func(t time.Time) {
			fmt.Println(t, true)
			time.Sleep(1)
		}})
	if err != nil {
		panic(err)
	}

	<-time.After(time.Minute)
	fmt.Printf("\n\ttest complete; stopping ticker...\n")

	ticker1.Stop()
	ticker2.Stop()
}
