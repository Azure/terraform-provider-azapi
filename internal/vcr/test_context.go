package vcr

// Maps goroutines to their active tests so the VCR recorder always knows which
// test context it is in when helpers such as BuildTestClient are called from a
// check, mostly to keep things correct when tests run in parallel.

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var goroutineTests sync.Map

// RegisterTestT associates t with the current goroutine.
// Called at the entry point of each VCR test run.
// Note: ensure a defer vcr.UnregisterTestT() is called immediately after this.
func RegisterTestT(t *testing.T) {
	goroutineTests.Store(currentGoroutineID(), t)
}

// UnregisterTestT removes the association for the current goroutine.
// Note: call in a defer immediately after RegisterTestT().
func UnregisterTestT() {
	goroutineTests.Delete(currentGoroutineID())
}

// CurrentTestName returns the name of the test running in this goroutine, or ""
// if none is registered.
func CurrentTestName() string {
	if v, ok := goroutineTests.Load(currentGoroutineID()); ok {
		return v.(*testing.T).Name()
	}
	return ""
}

func currentGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))
	id, _ := strconv.ParseInt(fields[0], 10, 64)
	return id
}
