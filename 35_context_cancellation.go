package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	// ---------- CONTEXT: WHAT IT'S FOR ----------

	//context.Context carries a cancellation signal (and optionally a
	//deadline, or request-scoped values) across function calls and
	//goroutines - it's how Go tells a whole tree of work "stop now"
	//without every function needing its own custom stop-channel.

	// ---------- WithCancel: MANUAL CANCELLATION ----------

	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx, "manual")
	time.Sleep(30 * time.Millisecond)
	cancel()                          // signals every goroutine using ctx to stop
	time.Sleep(20 * time.Millisecond) // let the worker's print land before main moves on

	// ---------- WithTimeout: AUTOMATIC CANCELLATION AFTER A DURATION ----------

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancelTimeout() // always call cancel, even on the timeout path - frees resources

	go worker(ctxTimeout, "timeout")
	time.Sleep(80 * time.Millisecond) // long enough to see the timeout fire

	// ---------- ctx.Err() SAYS WHY IT STOPPED ----------

	fmt.Println("manual ctx error:", ctx.Err())         // context canceled
	fmt.Println("timeout ctx error:", ctxTimeout.Err()) // context deadline exceeded

	// ---------- WithDeadline: CANCEL AT A SPECIFIC TIME ----------

	//same idea as WithTimeout, but you give an absolute time.Time instead
	//of a duration - useful when you're already tracking a deadline
	deadline := time.Now().Add(40 * time.Millisecond)
	ctxDeadline, cancelDeadline := context.WithDeadline(context.Background(), deadline)
	defer cancelDeadline()
	go worker(ctxDeadline, "deadline")
	time.Sleep(60 * time.Millisecond)

	// ---------- CANCELLATION PROPAGATES DOWN THE TREE ----------

	//cancelling a parent context automatically cancels every context
	//derived from it - you don't need to cancel each child by hand
	parent, cancelParent := context.WithCancel(context.Background())
	child, cancelChild := context.WithCancel(parent)
	defer cancelChild() // still good practice even though the parent will cancel it too

	go func() {
		<-child.Done()
		fmt.Println("child stopped because:", child.Err())
	}()
	cancelParent()                     // cancelling the parent cancels the child too
	time.Sleep(10 * time.Millisecond)

	// ---------- WithValue: REQUEST-SCOPED DATA (USE SPARINGLY) ----------

	//context.WithValue attaches a key/value pair to the context - meant for
	//cross-cutting request data like a request ID or auth token, NOT as a
	//general way to pass function arguments (that makes code harder to
	//follow and isn't type-checked by the compiler)
	type ctxKey string
	requestCtx := context.WithValue(context.Background(), ctxKey("requestID"), "abc-123")
	fmt.Println("requestID from context:", requestCtx.Value(ctxKey("requestID")))

	// ---------- GOTCHA: ALWAYS CALL cancel(), EVEN IF THE CONTEXT "FINISHED" ----------

	//WithCancel/WithTimeout/WithDeadline all return a cancel func - not
	//calling it leaks resources (an internal timer, in the timeout/deadline
	//case) until the parent context itself is done. `defer cancel()` right
	//after creating the context is the standard, safe habit.
	fmt.Println("done")
}

// worker keeps "working" until ctx says to stop - the standard shape for
// any function that should respect cancellation
func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done(): // closes when the context is canceled or times out
			fmt.Println(name, "worker stopping:", ctx.Err())
			return
		case <-time.After(10 * time.Millisecond):
			fmt.Println(name, "worker tick")
		}
	}
}
