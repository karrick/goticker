# goticker

Tiny Golang ticker library

## Description

A Ticker periodically invokes a callback function with the value of
the current time. Allows callers to optionally specify whether
invocations should occur at times that are rounded to the nearest
duration interval. A Ticker will continue until its Stop method is
invoked.

## Examples

### Some things need to happen periodically, but not on a specific rounded time:

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

### Some things need to happen on intervals rounded to nearest duration:

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
