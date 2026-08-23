package raceexample

import (
	"sync"
	"testing"
)

// under a normal `go test`, this test usually just PASSES - the race is
// still there, but nothing forces it to be OBSERVED. that's exactly the
// danger: a race condition can hide in a codebase for a long time,
// working "by luck" on every machine that happens to run it, until it
// doesn't. run this file with `go test -race` to have Go's runtime watch
// every memory access and prove the race exists, regardless of whether the
// final count happens to come out right this time.
func TestUnsafeCounter(t *testing.T) {
	c := &UnsafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()

	if c.Value() != 100 {
		t.Errorf("Value() = %d, want 100 (lost updates from the race)", c.Value())
	}
}

// the same test, against the Mutex-protected counter - this one should
// pass cleanly even under -race, since every access is properly synchronized.
func TestSafeCounter(t *testing.T) {
	c := &SafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()

	if c.Value() != 100 {
		t.Errorf("Value() = %d, want 100", c.Value())
	}
}
