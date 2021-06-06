# goticker

Tiny Golang ticker library

[![GoDoc](https://godoc.org/github.com/karrick/goticker?status.svg)](https://godoc.org/github.com/karrick/goticker)

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
    Duration: time.Minute,
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
// Rotate logs every midnight...
logTicker, err := goticker.New(goticker.Config{
    Round:    true,
    Duration: 24 * time.Hour,
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
creating in my programs. This certainly might not be remotely similar
to your boilerplate for similar purposes, but maybe it suits your
style as well.

When I want to create a task that needs to be run periodically with
the standard library, I always end up spawning a goroutine that loops
forever selecting on a time.Ticker's channel and yet another channel I
need to create to tell that goroutine when to stop the time.Ticker and
exit the goroutine.

This library serves as a place to hold the ticker boilerplate that I
tend to lean on using.

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
