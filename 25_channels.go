package main

import (
	"fmt"
	"sync"
	"time"
)

// takes a channel as a normal parameter - any function can send/receive
// on a channel it's given, which is how goroutines hand data back to main
func square(n int, result chan int) {
	result <- n * n
}

func main() {

	// ---------- CHANNELS: WHAT THEY ARE ----------

	//a channel is a typed pipe for sending values between goroutines.
	//make(chan T) creates one; its zero value (declared but not made) is nil.
	ch := make(chan int)

	go func() {
		ch <- 42 // send a value into the channel
	}()
	value := <-ch // receive a value from the channel
	fmt.Println("received:", value)

	// ---------- BLOCKING: SEND AND RECEIVE BOTH WAIT ----------

	//an unbuffered channel (no capacity) only completes a send once someone
	//is receiving at that exact moment, and vice versa - this is called a
	//rendezvous, and it's what makes plain channels double as a sync point.
	sync1 := make(chan string)
	go func() {
		fmt.Println("goroutine: about to send")
		sync1 <- "done"
		fmt.Println("goroutine: send completed - someone received it")
	}()
	time.Sleep(50 * time.Millisecond) // just so the prints above land first, for the demo
	fmt.Println("main: about to receive")
	fmt.Println("main: received", <-sync1)

	// ---------- PASSING CHANNELS INTO FUNCTIONS ----------

	resultCh := make(chan int)
	go square(6, resultCh)
	fmt.Println("square(6) =", <-resultCh)

	// ---------- CLOSING A CHANNEL ----------

	//close(ch) signals "nothing more is coming". only the SENDER should
	//close a channel - closing one you might still send on will panic later.
	nums := make(chan int)
	go func() {
		for i := 1; i <= 3; i++ {
			nums <- i
		}
		close(nums)
	}()

	//after close, receives still work but return immediately: real buffered
	//values first, then the zero value forever after that.
	fmt.Println("received:", <-nums, <-nums, <-nums)

	//the two-value receive form tells you whether that value was real (ok)
	//or the channel is closed and empty (!ok) - the only reliable way to
	//tell "closed" apart from "a real zero value was sent"
	v, ok := <-nums
	fmt.Println("after drained + closed - value:", v, "ok:", ok)

	// ---------- RANGE OVER A CHANNEL ----------

	//range receives values until the channel closes, then stops on its own
	letters := make(chan string)
	go func() {
		for _, l := range []string{"a", "b", "c"} {
			letters <- l
		}
		close(letters)
	}()
	fmt.Println("ranging over channel:")
	for l := range letters {
		fmt.Println(" ", l)
	}

	// ---------- MULTIPLE SENDERS, ONE CHANNEL (FAN-IN) ----------

	//several goroutines can safely send into the SAME channel - Go
	//guarantees each send/receive is handled as one atomic handoff
	results := make(chan int)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			results <- n * 10
		}(i)
	}

	//closing needs its own goroutine here: close() must happen only after
	//every sender is done, but main also needs to be ranging/receiving
	//concurrently, or the senders above would block forever with no one
	//reading yet
	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("fan-in results:")
	for r := range results {
		fmt.Println(" ", r)
	}

	// ---------- GOTCHAS ----------

	//1. sending on a CLOSED channel panics
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered:", r)
			}
		}()
		closedCh := make(chan int)
		close(closedCh)
		closedCh <- 1 // panic: send on closed channel
	}()

	//2. closing an ALREADY-CLOSED channel panics too
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered:", r)
			}
		}()
		closedCh := make(chan int)
		close(closedCh)
		close(closedCh) // panic: close of closed channel
	}()

	//3. a NIL channel (declared, never make()'d) blocks forever on both
	//send and receive - not an error, just permanently stuck. we launch
	//this in a goroutine and never wait for it, so it just quietly leaks
	//instead of freezing our program - don't do this on purpose in real code.
	var nilChan chan int
	go func() {
		<-nilChan // blocks here forever
	}()
	fmt.Println("(a goroutine is now blocked forever on a nil channel - harmless for this demo)")

	//4. DEADLOCK - not shown running here on purpose, since it crashes the
	//whole program instead of just panicking one goroutine:
	//   stuck := make(chan int)
	//   stuck <- 1 // no other goroutine will ever receive this
	//   -> fatal error: all goroutines are asleep - deadlock!
	//this happens whenever every goroutine is blocked waiting on a channel
	//and nothing is left running to unblock any of them.

	fmt.Println("done")
}
