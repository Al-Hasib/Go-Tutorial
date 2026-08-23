package main

import "fmt"

func main() {

	// ---------- BUFFERED CHANNELS: THE BASICS ----------

	//make(chan T, n) creates a channel with room for n values inside it.
	//a send only blocks once the buffer is full - unlike an unbuffered
	//channel, where every send blocks until someone receives.
	buffered := make(chan int, 3)

	buffered <- 1 // doesn't block - buffer has room
	buffered <- 2
	buffered <- 3
	fmt.Println("filled buffer without any receiver waiting")

	//len() = how many values are sitting in the buffer right now
	//cap() = the buffer's total capacity, fixed at creation
	fmt.Println("len:", len(buffered), "cap:", cap(buffered))

	fmt.Println("received:", <-buffered, <-buffered, <-buffered)
	fmt.Println("len after draining:", len(buffered))

	// ---------- WHEN A BUFFERED CHANNEL STILL BLOCKS ----------

	small := make(chan int, 2)
	small <- 1
	small <- 2
	fmt.Println("buffer full, len/cap:", len(small), "/", cap(small))

	//a 3rd send here would block forever (nobody is receiving, buffer is
	//full) - not run on purpose, since it would deadlock this program:
	//   small <- 3 // fatal error: all goroutines are asleep - deadlock!
	//receiving first frees up a slot, which is why buffered channels still
	//need a receiver eventually - they only delay blocking, they don't remove it
	fmt.Println("received one:", <-small)
	small <- 3 // now this fits, since draining one made room
	fmt.Println("received:", <-small, <-small)

	// ---------- BUFFERED VS UNBUFFERED: WHY IT MATTERS ----------

	//with make(chan int) (unbuffered), the line below would need a
	//goroutine ready to receive at the same instant, or it blocks.
	//with a buffer, a goroutine can fire off several sends and move on,
	//as long as it doesn't send more than the buffer can hold before
	//anyone catches up on receiving.
	queue := make(chan string, 5)
	tasks := []string{"email", "resize-image", "backup", "cleanup"}
	for _, t := range tasks {
		queue <- t // all 4 fit in a 5-slot buffer, so none of this blocks
	}
	close(queue) // safe to close now - every task is already queued

	fmt.Println("processing queued tasks:")
	for t := range queue {
		fmt.Println(" - processing:", t)
	}

	// ---------- A CLOSED BUFFERED CHANNEL STILL DRAINS ITS LEFTOVERS ----------

	//closing doesn't throw away what's already buffered - range/receive
	//still returns every real value first, THEN the zero value/!ok
	leftover := make(chan int, 3)
	leftover <- 100
	leftover <- 200
	close(leftover)

	for v := range leftover {
		fmt.Println("leftover value:", v) // both 100 and 200 come through
	}
	v, ok := <-leftover
	fmt.Println("after fully drained - value:", v, "ok:", ok)
}
