package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {

	// ---------- WHAT A RACE CONDITION IS ----------

	//a race condition happens when 2+ goroutines read/write the same
	//variable at the same time with no synchronization - the outcome
	//depends on timing, which makes it unpredictable and hard to debug.
	unsafeCounter := 0
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsafeCounter++ // NOT atomic: read, increment, write - 3 separate steps
		}()
	}
	wg.Wait()

	//this often prints something less than 1000, because two goroutines
	//can both read the same old value before either writes it back,
	//silently losing an increment. it won't happen on every single run -
	//that unpredictability IS the problem. to prove it reliably, run:
	//   go run -race 31_race_conditions.go
	fmt.Println("unsafe counter (expected 1000, may be less):", unsafeCounter)

	// ---------- FIX 1: sync.Mutex ----------

	safeCounter := 0
	var mu sync.Mutex
	var wg2 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			mu.Lock()
			safeCounter++
			mu.Unlock()
		}()
	}
	wg2.Wait()
	fmt.Println("mutex-protected counter:", safeCounter)

	// ---------- FIX 2: sync/atomic - lighter weight for simple cases ----------

	//atomic operations do the read-modify-write as one indivisible CPU-level
	//step - no Lock/Unlock needed, but this only works for simple
	//counters/flags, not for protecting a whole block of logic or several
	//variables together
	var atomicCounter int64
	var wg3 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			atomic.AddInt64(&atomicCounter, 1)
		}()
	}
	wg3.Wait()
	fmt.Println("atomic counter:", atomic.LoadInt64(&atomicCounter))

	// ---------- DETECTING RACES: THE -race FLAG ----------

	//`go run -race file.go`, `go test -race`, or `go build -race` compile in
	//instrumentation that watches every memory access and reports exactly
	//which goroutines raced on which variable, with stack traces. it's
	//slower, so it's used in testing/CI rather than production - but it's
	//the standard way to actually confirm a race, instead of guessing from
	//inconsistent output.
	fmt.Println("try: go run -race 31_race_conditions.go")
}
