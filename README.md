# goticker

Tiny Golang ticker library

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GoDoc](https://godoc.org/github.com/karrick/goticker?status.svg)](https://godoc.org/github.com/karrick/goticker)
[![GoReportCard](https://goreportcard.com/badge/github.com/karrick/goticker)](https://goreportcard.com/report/github.com/karrick/goticker)

## Description

A Ticker periodically invokes a callback function with the value of
the current time. Allows callers to optionally specify whether
invocations should occur at times that are rounded to the nearest
duration interval. A Ticker will continue until its Stop method is
invoked.

## Examples

### When things need to happen periodically, but not on a specific rounded time:

```Go
// Emit metrics every minute...
metricTicker, err := goticker.New(goticker.Config{
    Interval: time.Minute,
    Callback: func(t time.Time) {
        metrics.Emit()
    }})
if err != nil {
    panic(err) // TODO: handle appropriately
}

// some time later...
metricTicker.Stop()
```

### When things need to happen on intervals rounded to nearest duration:

```Go
// Rotate logs every 24 hours, on the hour...
logTicker, err := goticker.New(goticker.Config{
    Round:    time.Hour,
    Interval: 24 * time.Hour,
    Callback: func(t time.Time) {
        logger.Rotate()
    }})
if err != nil {
    panic(err) // TODO: handle appropriately
}

// some time later...
logTicker.Stop()
```

A slightly different variation:

```Go
// Rotate logs every midnight...
logTicker, err := goticker.New(goticker.Config{
    Round:    24 * time.Hour,
    Interval: 24 * time.Hour,
    Callback: func(t time.Time) {
        logger.Rotate()
    }})
if err != nil {
    panic(err) // TODO: handle appropriately
}

// some time later...
logTicker.Stop()
```

## Why?

I created this library to eliminate boilerplate code I always end up
creating in programs. This certainly might not be similar to your
boilerplate for similar purposes, but maybe it suits your style as
well, or maybe it just serves an inspiration to a much better way to
manage some sort of event ticker.

When I want to create a task that needs to be run periodically with
the standard library, I always end up spawning a goroutine that loops
forever selecting on a time.Ticker's channel and yet another channel I
need to create to tell that goroutine when to stop the time.Ticker and
exit the goroutine. This boilerplate is needed because time.Ticker
provides a channel to receive ticker events from, rather than invoking
a provided function. See below for an example of the boilerplate I
wanted to eliminate.

```Go
func NewTicker(duration time.Duration, callback func(time.Time)) chan struct{} {
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

func main() {
    cancel := NewTicker(time.Second, func(t time.Time) {
        fmt.Print(".")
        hr, min, sec := t.Clock()
        if sec == 0 && min == 0 && hr == 0 {
            fmt.Printf("\n\tMidnight %v\n, t)
        }
    })

    <-time.After(2 * time.Minute)
    fmt.Printf("\n\ttest complete; stopping ticker...\n")
    close(cancel)
}
```

While this library originally served as a place to hold the above
ticker boilerplate that I tend on using, it has evolved to be a tad
more flexible. Rather than create boilerplate that loops over the Go
standard library time.Ticker and a few channels, this library does not
use time.Ticker, and for that reason does not need to use
channels. Instead, this library allows the caller to specify a
function to invoke when the duration time has elapsed, and calls it
after sleeping for the required time.
