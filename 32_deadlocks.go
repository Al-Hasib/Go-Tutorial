package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	// ---------- WHAT A DEADLOCK IS ----------

	//a deadlock is everyone waiting on everyone else, forever - no
	//goroutine can make progress because each one needs something that
	//only another blocked goroutine could provide.

	// ---------- CASE 1: TWO LOCKS TAKEN IN OPPOSITE ORDER ----------

	//goroutine A locks muA then wants muB; goroutine B locks muB then wants
	//muA - if both grab their first lock before either reaches the second,
	//neither can ever proceed.
	//
	//running this for real would hang the program, so it's wrapped in a
	//timeout: if it doesn't finish quickly, we know it deadlocked, and we
	//move on (the 2 stuck goroutines just leak - fine for a demo, never
	//for real code).
	var muA, muB sync.Mutex
	deadlockDone := make(chan bool)

	go func() {
		muA.Lock()
		time.Sleep(20 * time.Millisecond) // give goroutine B time to lock muB first
		muB.Lock()
		muB.Unlock()
		muA.Unlock()
		deadlockDone <- true
	}()
	go func() {
		muB.Lock()
		time.Sleep(20 * time.Millisecond) // give goroutine A time to lock muA first
		muA.Lock()
		muA.Unlock()
		muB.Unlock()
		deadlockDone <- true
	}()

	select {
	case <-deadlockDone:
		fmt.Println("finished without deadlocking this time")
	case <-time.After(200 * time.Millisecond):
		fmt.Println("deadlock detected: both goroutines are stuck waiting on each other's lock")
	}

	// ---------- THE FIX: ALWAYS LOCK IN THE SAME ORDER ----------

	//NOTE: using fresh mutexes here, not muA/muB - those two are left
	//locked forever by the leaked goroutines above, so reusing them would
	//deadlock this section for real instead of just demonstrating the fix.
	//
	//if EVERY goroutine locks muC before muD, this exact deadlock can't
	//happen - whoever gets muC first is guaranteed to also get muD, since
	//nobody else will be holding muD while waiting for muC.
	var muC, muD sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			muC.Lock()
			muD.Lock()
			muD.Unlock()
			muC.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("consistent lock ordering finished cleanly")

	// ---------- CASE 2: A CHANNEL WITH NO ONE TO RECEIVE ----------

	//unlike case 1, Go's runtime can actually detect this one: if EVERY
	//goroutine in the whole program is blocked (not just some of them), Go
	//notices nothing could ever make progress and kills the program
	//immediately instead of hanging forever silently:
	//   stuck := make(chan int)
	//   stuck <- 1 // fatal error: all goroutines are asleep - deadlock!
	//not run here on purpose, since it would crash this whole file.

	// ---------- WHY CASE 1 DIDN'T CRASH BUT CASE 2 WOULD ----------

	//in case 1, main and other goroutines were still doing work (or could),
	//so Go's runtime saw the program wasn't FULLY stuck - only those two
	//goroutines were, so they just leaked silently. Go only auto-detects
	//and crashes when literally nothing in the whole process can proceed.
	fmt.Println("done")
}
