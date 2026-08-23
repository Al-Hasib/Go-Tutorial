package main

import (
	"fmt"
	"sync"
)

func main() {

	// ---------- MUTEX: PROTECTING SHARED STATE ----------

	//sync.Mutex's zero value is ready to use - no constructor needed
	var mu sync.Mutex
	counter := 0

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++ // safe: only one goroutine can be inside Lock/Unlock at a time
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("counter after 1000 concurrent increments:", counter)

	// ---------- defer mu.Unlock() IS THE SAFER HABIT ----------

	safeIncrement := func(mu *sync.Mutex, n *int) {
		mu.Lock()
		defer mu.Unlock() // still runs even if this function panics partway through
		*n++
	}
	var mu2 sync.Mutex
	total := 0
	var wg2 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			safeIncrement(&mu2, &total)
		}()
	}
	wg2.Wait()
	fmt.Println("total via safeIncrement:", total)

	// ---------- PROTECTING A MAP ----------

	//maps are NOT safe for concurrent read/write - a mutex fixes that too
	var mapMu sync.Mutex
	visits := make(map[string]int)

	pages := []string{"home", "about", "home", "contact", "home", "about"}
	var wg3 sync.WaitGroup
	for _, page := range pages {
		wg3.Add(1)
		go func(p string) {
			defer wg3.Done()
			mapMu.Lock()
			visits[p]++
			mapMu.Unlock()
		}(page)
	}
	wg3.Wait()
	fmt.Println("visits:", visits)

	// ---------- TryLock: CHECK WITHOUT BLOCKING ----------

	var mu3 sync.Mutex
	mu3.Lock()
	if mu3.TryLock() {
		fmt.Println("acquired lock (won't happen, it's already locked)")
	} else {
		fmt.Println("TryLock failed instantly instead of blocking - already locked")
	}
	mu3.Unlock()

	// ---------- GOTCHA: FORGETTING Unlock() DEADLOCKS EVERYONE ----------

	//a Mutex left locked forever blocks every future Lock() call on it,
	//including from the SAME goroutine - it is not reentrant. not run here,
	//since it would hang this lesson forever:
	//   var stuck sync.Mutex
	//   stuck.Lock()
	//   stuck.Lock() // blocks forever
	fmt.Println("done")
}
