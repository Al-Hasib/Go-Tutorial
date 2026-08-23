package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const dsn = "postgres://postgres:postgres@localhost:55432/tutorial?sslmode=disable"

func main() {

	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	// ---------- *sql.DB IS ALREADY A POOL ----------

	//sql.Open doesn't hand you one connection - it hands you a POOL,
	//already safe for many goroutines to use concurrently. real
	//connections open lazily, as concurrent work actually needs them.

	// ---------- TUNING THE POOL ----------

	db.SetMaxOpenConns(3)                  // never more than 3 real connections to Postgres at once
	db.SetMaxIdleConns(3)                  // keep up to 3 idle ones ready, instead of reconnecting each time
	db.SetConnMaxLifetime(5 * time.Minute) // force-recycle a connection after this long
	db.SetConnMaxIdleTime(1 * time.Minute) // close an idle connection that's sat unused this long

	// ---------- WATCHING THE POOL UNDER LOAD ----------

	//pg_sleep(0.3) simulates a slow query - firing 6 of these at once
	//against a 3-connection pool means half of them have to wait for a
	//connection to free up
	var wg sync.WaitGroup
	start := time.Now()
	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if _, err := db.Exec(`SELECT pg_sleep(0.3)`); err != nil {
				fmt.Println("query", n, "error:", err)
				return
			}
			fmt.Printf("query %d finished at +%s\n", n, time.Since(start).Round(10*time.Millisecond))
		}(i)
	}

	//poll pool stats while those 6 queries are in flight
	for i := 0; i < 4; i++ {
		time.Sleep(150 * time.Millisecond)
		s := db.Stats()
		fmt.Printf("  [stats] open=%d inUse=%d idle=%d waitCount=%d\n", s.OpenConnections, s.InUse, s.Idle, s.WaitCount)
	}

	wg.Wait()
	fmt.Println("\ntotal time for 6 queries over a 3-connection pool:", time.Since(start).Round(10*time.Millisecond))
	//if pooling worked as expected, this lands near 2 batches of ~0.3s
	//(~0.6s total) - not 6 sequential queries (~1.8s, no pooling benefit)
	//and not all 6 running fully in parallel (~0.3s, which would need 6
	//connections) - MaxOpenConns=3 explains that middle number.
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
