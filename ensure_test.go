package goticker

import (
	"strings"
	"testing"
)

func ensurePanic(tb testing.TB, label, want string, callback func()) {
	tb.Helper()
	defer func() {
		r := recover()
		if r == nil || !strings.Contains(r.(error).Error(), want) {
			tb.Errorf("TEST: %s; GOT: %v; WANT: %v", label, r, want)
		}
	}()
	callback()
}

// ensureNoPanic prettifies the output so one knows which test case caused a
// panic.
func ensureNoPanic(tb testing.TB, label string, callback func()) {
	tb.Helper()
	defer func() {
		if r := recover(); r != nil {
			tb.Fatalf("TEST: %s: GOT: %v", label, r)
		}
	}()
	callback()
}
