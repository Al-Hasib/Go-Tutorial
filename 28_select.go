package main

import (
	"fmt"
	"time"
)

func main() {

	// ---------- SELECT: THE BASICS ----------

	chA := make(chan string)
	chB := make(chan string)

	go func() {
		time.Sleep(30 * time.Millisecond)
		chA <- "from A"
	}()
	go func() {
		time.Sleep(10 * time.Millisecond)
		chB <- "from B"
	}()

	//select waits on multiple channel operations at once and runs whichever
	//case becomes ready first - like a switch, but for channels
	for i := 0; i < 2; i++ {
		select {
		case a := <-chA:
			fmt.Println("select got:", a)
		case b := <-chB:
			fmt.Println("select got:", b)
		}
	}

	// ---------- WHEN MULTIPLE CASES ARE READY AT ONCE ----------

	bothReady := make(chan int, 1)
	alsoReady := make(chan int, 1)
	bothReady <- 1
	alsoReady <- 2

	//both cases are ready immediately here - select picks one at RANDOM,
	//not top-to-bottom like a normal switch. run this a few times and the
	//printed case can differ
	select {
	case v := <-bothReady:
		fmt.Println("picked bothReady:", v)
	case v := <-alsoReady:
		fmt.Println("picked alsoReady:", v)
	}

	// ---------- default: NON-BLOCKING SELECT ----------

	maybe := make(chan int)
	select {
	case v := <-maybe:
		fmt.Println("got:", v)
	default:
		//default runs immediately if no other case is ready - this is what
		//makes select never block
		fmt.Println("nothing ready, moving on")
	}

	// ---------- TIMEOUTS WITH time.After ----------

	slow := make(chan string)
	select {
	case r := <-slow:
		fmt.Println("got:", r)
	case <-time.After(50 * time.Millisecond):
		//time.After returns a channel that fires once, after the given
		//duration - a clean way to bound how long you'll wait
		fmt.Println("timed out waiting for slow")
	}

	// ---------- SELECT CAN SEND TOO, NOT JUST RECEIVE ----------

	outbox := make(chan int, 1)
	select {
	case outbox <- 99:
		fmt.Println("sent 99 into outbox")
	default:
		fmt.Println("outbox full, skipped send")
	}

	// ---------- FOR-SELECT: LISTENING UNTIL TOLD TO STOP ----------

	done := make(chan struct{})
	ticks := make(chan int)

	go func() {
		for i := 1; i <= 3; i++ {
			ticks <- i
			time.Sleep(10 * time.Millisecond)
		}
		close(done) // signal "stop now"
	}()

	fmt.Println("for-select loop:")
loop:
	for {
		select {
		case t := <-ticks:
			fmt.Println("  tick:", t)
		case <-done:
			fmt.Println("  told to stop")
			break loop
		}
	}

	// ---------- GOTCHAS ----------

	//an empty select{} blocks forever - sometimes used deliberately to
	//park a program, but usually a mistake. not run here since it would
	//hang this lesson:
	//   select {}

	//a nil channel case in a select is never chosen (nil channels block
	//forever) - this is actually useful: set a case's channel to nil to
	//dynamically "turn it off" without removing the case
	var disabled chan int // nil
	select {
	case <-disabled:
		fmt.Println("this can never print")
	default:
		fmt.Println("disabled case was skipped, as expected")
	}
}
