// Package raceexample holds two counters side by side - one with a real
// data race, one without - so the Race Detector lesson can show the
// difference on real code instead of just describing it.
package raceexample

import "sync"

// UnsafeCounter increments its value with no synchronization at all.
type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Increment() { c.value++ }
func (c *UnsafeCounter) Value() int { return c.value }

// SafeCounter protects the same value with a Mutex (see the Mutex lesson
// in Phase 5 for more on this).
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}
