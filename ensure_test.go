package goticker

import (
	"strings"
	"testing"
)

func ensureError(tb testing.TB, err error, contains ...string) {
	tb.Helper()
	if len(contains) == 0 || (len(contains) == 1 && contains[0] == "") {
		if err != nil {
			tb.Fatalf("GOT: %v; WANT: %v", err, contains)
		}
	} else if err == nil {
		tb.Errorf("GOT: %v; WANT: %v", err, contains)
	} else {
		for _, stub := range contains {
			if stub != "" && !strings.Contains(err.Error(), stub) {
				tb.Errorf("GOT: %v; WANT: %q", err, stub)
			}
		}
	}
}

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
