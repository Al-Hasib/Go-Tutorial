package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // an in-flight request the shutdown below has to wait on
		fmt.Fprintln(w, "finally done")
	})

	server := &http.Server{Handler: mux}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	must(err)
	addr := "http://" + listener.Addr().String()

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error:", err)
		}
	}()
	fmt.Println("server listening on", addr)

	// ---------- LISTENING FOR THE SIGNALS THAT MEAN "STOP" ----------

	//signal.NotifyContext returns a context that's cancelled the moment the
	//process receives any of the given signals - Ctrl+C sends os.Interrupt;
	//SIGTERM is what most deployment platforms (Docker, Kubernetes, systemd)
	//send to ask a process to shut down cleanly before killing it outright
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//everything below this line, up to demoCtx, is exactly what real
	//production code looks like: block until a real signal arrives, then
	//shut down. the demoCtx wrapper exists ONLY so this lesson finishes on
	//its own when run non-interactively - delete it and this is production code.
	demoCtx, cancelDemo := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancelDemo()

	// ---------- SIMULATING TRAFFIC WHILE WE WAIT ----------

	go func() {
		resp, err := http.Get(addr + "/slow")
		if err == nil {
			fmt.Println("in-flight request finished with status:", resp.StatusCode)
			resp.Body.Close()
		}
	}()
	time.Sleep(20 * time.Millisecond) // let that request actually start before shutdown begins

	<-demoCtx.Done()
	fmt.Println("\nshutdown triggered - in real usage this line only runs after a real OS signal")

	// ---------- THE ACTUAL SHUTDOWN SEQUENCE ----------

	//give in-flight requests (like /slow above) a bounded amount of time to
	//finish naturally, instead of either waiting forever or cutting them off
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	fmt.Println("closing new connections, waiting for in-flight requests...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("graceful shutdown failed, forcing close:", err)
		server.Close() // last resort: drop everything immediately
	} else {
		fmt.Println("shut down cleanly")
	}

	// ---------- WHAT A REAL PROGRAM DOES AFTER THIS ----------

	//close a database pool (db.Close()), flush any buffered logs/metrics,
	//release a distributed lock, etc - anything else the process was
	//holding onto also belongs after Shutdown() returns, in this same
	//"cleaning up before exit" phase.
	fmt.Println("done")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
