package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func main() {

	// ---------- BUILDING A REAL HTTP SERVER ----------

	//an http.ServeMux is a request router: it matches an incoming request's
	//method+path against registered patterns and calls the matching handler
	mux := http.NewServeMux()

	//"/{$}" matches ONLY the exact root path - a plain "/" would match the
	//root AND every unmatched path underneath it (a catch-all), which would
	//swallow the "unregistered path" demo below instead of letting it 404
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "welcome to the home page")
	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		//http.ResponseWriter is how a handler sends a response back:
		//Header() to set response headers, WriteHeader() for the status
		//code, Write()/fmt.Fprint* for the body
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "hello, HTTP!")
	})

	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // simulates slow work, used by the shutdown demo below
		fmt.Fprintln(w, "finally done")
	})

	//http.Server bundles the router with real production-relevant settings -
	//timeouts matter a lot: without them, a slow or malicious client can tie
	//up a connection (and a goroutine) indefinitely
	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//net.Listen on port 0 asks the OS for any free port - avoids hardcoding
	//a port that might already be busy, which is what lets this lesson run
	//repeatedly without conflicts
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	addr := "http://" + listener.Addr().String()
	fmt.Println("server listening on", addr)

	//ListenAndServe (or Serve, given our own listener) BLOCKS the calling
	//goroutine forever - it only returns on an error or after Shutdown() is
	//called. that's why it always runs in its own goroutine, leaving main
	//free to keep going (here: to act as the client hitting it, then to
	//shut it down).
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error:", err)
		}
	}()

	// ---------- TALKING TO OUR OWN SERVER ----------

	resp, _ := http.Get(addr + "/")
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Print("GET / -> ", string(body))

	resp, _ = http.Get(addr + "/hello")
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Print("GET /hello -> ", string(body))

	//an unregistered path falls through to ServeMux's built-in 404
	resp, _ = http.Get(addr + "/missing")
	fmt.Println("GET /missing -> status:", resp.StatusCode)
	resp.Body.Close()

	// ---------- GRACEFUL SHUTDOWN ----------

	//kick off a slow request in the background, then ask the server to shut
	//down right after - Shutdown() waits for in-flight requests (like this
	//one) to finish instead of just killing every connection outright
	go func() {
		resp, err := http.Get(addr + "/slow")
		if err == nil {
			fmt.Println("slow request finished with status:", resp.StatusCode)
			resp.Body.Close()
		}
	}()
	time.Sleep(20 * time.Millisecond) // let the /slow request actually start first

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	fmt.Println("shutting down...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("shutdown error:", err)
	} else {
		fmt.Println("server shut down cleanly (waited for the slow request)")
	}
}
