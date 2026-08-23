package main

import (
	"fmt"
	"sync"
)

// ---------- FAN-OUT: ONE SOURCE CHANNEL, MULTIPLE WORKER GOROUTINES ----------

// fanOut starts n workers that all read from the same "in" channel, so the
// work is spread across all of them concurrently instead of processed one
// item at a time.
func fanOut(in <-chan int, n int, work func(int) int) []<-chan int {
	outs := make([]<-chan int, n)
	for i := 0; i < n; i++ {
		out := make(chan int)
		outs[i] = out
		go func(out chan int) {
			defer close(out)
			for v := range in {
				out <- work(v)
			}
		}(out)
	}
	return outs
}

// ---------- FAN-IN: MULTIPLE CHANNELS MERGED INTO ONE ----------

// fanIn merges several input channels into a single output channel, so a
// consumer can read from one place instead of juggling many.
func fanIn(channels ...<-chan int) <-chan int {
	merged := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, c := range channels {
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				merged <- v
			}
		}(c)
	}

	//closer: only close merged once every input channel has been fully drained
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {

	// ---------- SOURCE ----------

	source := make(chan int)
	go func() {
		for i := 1; i <= 9; i++ {
			source <- i
		}
		close(source)
	}()

	// ---------- FAN-OUT: 3 WORKERS SHARE THE SOURCE ----------

	square := func(n int) int { return n * n }
	workerOutputs := fanOut(source, 3, square)
	fmt.Println("fanned out across", len(workerOutputs), "workers")

	// ---------- FAN-IN: MERGE ALL 3 WORKER OUTPUTS INTO ONE CHANNEL ----------

	merged := fanIn(workerOutputs...)

	sum := 0
	count := 0
	for v := range merged {
		sum += v
		count++
	}
	fmt.Println("received", count, "results, sum:", sum) // order varies, sum/count don't

	// ---------- WHY BOTHER: FAN-OUT/FAN-IN IN ONE PICTURE ----------

	// source -> [worker 1] -\
	//        -> [worker 2] --> fan-in -> single results channel -> consumer
	//        -> [worker 3] -/
	//
	//fan-out spreads CPU/IO-bound work across several goroutines; fan-in
	//gives the consumer one place to read results from, without caring how
	//many workers actually produced them.
}
