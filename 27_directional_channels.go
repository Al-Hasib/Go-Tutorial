package main

import "fmt"

// chan<- int : SEND-ONLY - this function may only send into out, never read from it
func producer(out chan<- int, count int) {
	for i := 1; i <= count; i++ {
		out <- i
	}
	close(out) // a send-only channel can still be closed by the sender
}

// <-chan int : RECEIVE-ONLY - this function may only read from in, never send
func consumer(in <-chan int) int {
	sum := 0
	for n := range in {
		sum += n
	}
	return sum
}

// a function can also RETURN a receive-only channel - the caller gets a
// value it can read from, but has no way to send on it or close it itself.
// this is a common pattern: the function owns the channel and controls it.
func generator(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 1; i <= count; i++ {
			out <- i * i
		}
		close(out)
	}()
	return out
}

func main() {

	// ---------- DIRECTIONAL CHANNELS: WHY THEY EXIST ----------

	//a plain "chan int" can both send and receive. narrowing a function
	//parameter to chan<- int (send-only) or <-chan int (receive-only)
	//documents what that function is allowed to do with it - and the
	//compiler enforces it, catching misuse before the program even runs.
	nums := make(chan int) // ordinary bidirectional channel

	go producer(nums, 3) // passed where a send-only parameter is expected -
	// Go converts chan int -> chan<- int automatically here, no cast needed

	total := consumer(nums) // passed where a receive-only parameter is expected -
	// same automatic chan int -> <-chan int conversion
	fmt.Println("sum from producer/consumer:", total)

	//what the compiler would refuse to build, if you tried it:
	//   func producer(out chan<- int, count int) { ...; x := <-out; ... }
	//   -> compile error: invalid operation: cannot receive from send-only channel out
	//that's the entire point - it turns a potential runtime bug into a
	//build-time error instead.

	// ---------- RETURNING A RECEIVE-ONLY CHANNEL ----------

	squares := generator(5) // caller only ever sees a <-chan int
	fmt.Println("squares from generator:")
	for sq := range squares {
		fmt.Println(" ", sq)
	}

	//the caller has no way to do "squares <- 99" or "close(squares)" here -
	//generator() alone owns and controls that channel's lifecycle, which
	//avoids a whole class of "who's responsible for closing this?" bugs

	// ---------- CHAINING: A SMALL PIPELINE ----------

	//each stage only receives from the previous one and sends to the next -
	//directional types make each stage's role obvious just from its signature
	pipelineOut := double(generator(4))
	fmt.Println("piped through double():")
	for v := range pipelineOut {
		fmt.Println(" ", v)
	}
}

// takes receive-only (reads from the previous stage), returns receive-only
// (the next stage, or main, can only read what this produces)
func double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			out <- v * 2
		}
		close(out)
	}()
	return out
}
