package main

import (
	"fmt"
	"sync"
)

func main() {

	// ---------- WORKER POOL: A FIXED NUMBER OF WORKERS SHARE ONE JOB QUEUE ----------

	//instead of one goroutine per job (which could mean thousands at once),
	//a worker pool caps concurrency: N workers pull from the same jobs
	//channel, so at most N jobs run at the same time.
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup

	const workerCount = 3
	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs) // no more jobs - lets each worker's "range jobs" loop end

	wg.Wait()      // wait for every worker to finish draining jobs
	close(results) // safe now - nothing will send to results anymore

	fmt.Println("worker pool results:")
	for r := range results {
		fmt.Println(" ", r)
	}

	// ---------- WHY A FIXED POOL, NOT ONE GOROUTINE PER JOB ----------

	//unbounded goroutines can exhaust memory or overwhelm a downstream
	//resource (a database, an API with a rate limit, disk I/O). a worker
	//pool caps how much runs at once, no matter how many jobs arrive.

	// ---------- A SIMPLER ALTERNATIVE: A SEMAPHORE ----------

	//for one-off bursts of work where you don't need persistent workers, a
	//buffered channel of empty structs can cap concurrency just as well
	sem := make(chan struct{}, 2) // only 2 "slots" available at once
	var wg2 sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg2.Add(1)
		go func(n int) {
			defer wg2.Done()
			sem <- struct{}{}        // take a slot - blocks if all slots are taken
			defer func() { <-sem }() // release the slot when done
			fmt.Println("task", n, "running (max 2 at a time)")
		}(i)
	}
	wg2.Wait()
}

// worker pulls jobs off the shared channel until it's closed and empty,
// squares each one, and sends the result onward
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		results <- j * j
	}
}
