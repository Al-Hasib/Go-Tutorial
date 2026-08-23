package main

import (
	"fmt"
	"sync"
	"time"
)

//a normal function - nothing special about it yet
func sayHello() {
	fmt.Println("hello from a goroutine")
}

func main() {

	// ---------- GOROUTINES: THE BASICS ----------

	//calling a function normally - Go runs it and waits for it to finish
	//before moving to the next line
	sayHello()

	//putting "go" in front of a function call runs it concurrently, in its
	//own goroutine, and returns immediately - main does NOT wait for it
	go sayHello()
	fmt.Println("this line can print before or after the goroutine above")

	//same thing works with an inline (anonymous) function too
	go func() {
		fmt.Println("hello from an anonymous goroutine")
	}()

	//PROBLEM: main is the program. the moment main() returns, the whole
	//program exits - even if some goroutines haven't run yet.
	//time.Sleep here just buys them enough time to run before that happens.
	//(this is NOT how you're supposed to wait for goroutines in real code -
	//sleeping a fixed amount is a guess, not a guarantee. the proper way,
	//sync.WaitGroup, is next lesson.)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("main is done")

	// ---------- LAUNCHING SEVERAL AT ONCE ----------

	//each loop iteration starts its own goroutine - they all run at the
	//same time, so the printed order is not guaranteed to be 1, 2, 3
	for i := 1; i <= 3; i++ {
		go fmt.Println("worker number", i)
	}

	time.Sleep(100 * time.Millisecond) // again, just to let them finish before we exit
	fmt.Println("all workers had a chance to run")

	// ---------- WAITGROUP: WAITING PROPERLY ----------

	//time.Sleep above was a guess - it works, but only because we guessed
	//a long enough delay. sync.WaitGroup lets goroutines report "I'm done"
	//so main can wait for an exact, real signal instead of guessing.
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // "one more goroutine to wait for" - call this BEFORE starting it

		//i is passed in as a parameter to be explicit about it. since Go 1.22,
		//each loop iteration already gets its own copy of i automatically, so
		//this exact bug is mostly historical - but you'll still see code
		//(and job interview questions) about it, since older Go shared one i
		//across every iteration, and every goroutine below would've printed
		//the same final value instead of 1, 2, 3.
		go func(n int) {
			defer wg.Done() // "this one is done" - defer makes sure it always runs
			fmt.Println("worker", n, "working")
		}(i)
	}

	wg.Wait() // blocks here until all 3 have called Done() - no guessing, no sleep
	fmt.Println("all workers finished (WaitGroup confirmed it, not a guess)")

	//note: a sync.WaitGroup must never be copied after it's used - if you
	//pass one into another function, pass a *pointer* (func(wg *sync.WaitGroup))
	//so every goroutine shares the exact same counter.

	// ---------- GOTCHA: A PANIC IN A GOROUTINE CAN KILL THE WHOLE PROGRAM ----------

	//goroutines are cheap - a few KB of stack that grows as needed, so
	//spawning thousands of them is normal (unlike OS threads, which are
	//much heavier). that cheapness is part of why Go code reaches for
	//goroutines so freely.
	var wg3 sync.WaitGroup
	wg3.Add(1)
	go func() {
		defer wg3.Done()

		//recover only protects the goroutine it's called in - it does NOT
		//catch panics from OTHER goroutines. if this goroutine panicked
		//without its own recover, the entire program would crash, even
		//though main and every other goroutine are perfectly fine.
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered inside the goroutine:", r)
			}
		}()

		panic("something went wrong in this goroutine")
	}()
	wg3.Wait()
	fmt.Println("main kept running because the goroutine recovered its own panic")
}
