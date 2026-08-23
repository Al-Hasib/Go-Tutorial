package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	// ---------- RWMUTEX: MANY READERS, ONE WRITER ----------

	//RWMutex separates "reading" locks from "writing" locks:
	// - RLock/RUnlock: any number of readers can hold this at once
	// - Lock/Unlock: a writer needs EXCLUSIVE access - blocks everyone else
	var mu sync.RWMutex
	data := map[string]int{"a": 1, "b": 2}

	var wg sync.WaitGroup

	//several concurrent readers - all can be inside RLock at the same time
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.RLock()
			defer mu.RUnlock()
			time.Sleep(10 * time.Millisecond) // hold the read lock briefly, on purpose
			fmt.Println("reader", id, "sees:", data)
		}(i)
	}

	//one writer - needs to wait for readers to finish, and blocks any new
	//readers until it's done
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond) // let readers start first, for the demo
		mu.Lock()
		defer mu.Unlock()
		data["c"] = 3
		fmt.Println("writer added c")
	}()

	wg.Wait()
	fmt.Println("final data:", data)

	// ---------- WHEN TO USE RWMutex OVER Mutex ----------

	//RWMutex only pays off when reads vastly outnumber writes - for
	//balanced or write-heavy workloads, a plain sync.Mutex is simpler and
	//often just as fast, since RWMutex carries more internal bookkeeping.

	// ---------- GOTCHA: UPGRADING RLock TO Lock DEADLOCKS ----------

	//calling Lock() while the SAME goroutine still holds an RLock() never
	//finishes - Go doesn't support "upgrading" a read lock to a write lock.
	//not run here, since it would hang forever:
	//   mu.RLock()
	//   mu.Lock() // deadlock: waiting for a write lock while holding a read lock
	//   mu.Unlock()
	//   mu.RUnlock()
	fmt.Println("done")
}
