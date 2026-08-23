package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

// Middleware wraps a handler with extra behavior, returning a new handler
// that runs that behavior before/after calling the original one.
type Middleware func(http.Handler) http.Handler

// ---------- LOGGING MIDDLEWARE ----------

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // call onward to whatever this wraps
		log.Printf("%s %s took %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// ---------- RECOVERY MIDDLEWARE ----------

// without this, a panic inside a handler could take down request handling
// entirely - writing your own recovery middleware makes the behavior
// explicit and lets you decide exactly what the client sees (a clean 500,
// not a stack trace).
func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("recovered from panic: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "something went wrong")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ---------- A MIDDLEWARE THAT ADDS A HEADER TO EVERY RESPONSE ----------

func withServerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "go-tutorial")
		next.ServeHTTP(w, r)
	})
}

// chain applies middlewares in the order listed - the first one wraps
// everything else, so it runs FIRST on the way in and LAST on the way out.
func chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello!")
	})

	mux.HandleFunc("GET /boom", func(w http.ResponseWriter, r *http.Request) {
		panic("something exploded in this handler")
	})

	// ---------- APPLYING MIDDLEWARE TO THE WHOLE ROUTER ----------

	//order matters: withRecovery is listed first, so it wraps everything
	//else and its deferred recover runs LAST on the way out - meaning it can
	//catch a panic from any handler, AND from withLogging/withServerHeader too
	handler := chain(mux, withRecovery, withLogging, withServerHeader)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, _ := http.Get(server.URL + "/hello")
	fmt.Println("GET /hello -> status:", resp.StatusCode, "X-Powered-By:", resp.Header.Get("X-Powered-By"))
	resp.Body.Close()

	resp, _ = http.Get(server.URL + "/boom")
	fmt.Println("GET /boom -> status:", resp.StatusCode, "(server is still alive, thanks to withRecovery)")
	resp.Body.Close()

	// ---------- PER-ROUTE MIDDLEWARE ----------

	//middleware doesn't have to apply globally - wrap just one handler to
	//scope it to a single route
	mux2 := http.NewServeMux()
	mux2.Handle("GET /admin", chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "admin area")
		}),
		withLogging, // only this route gets logged, in this second example
	))
	mux2.HandleFunc("GET /public", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "public area")
	})

	server2 := httptest.NewServer(mux2)
	defer server2.Close()
	resp, _ = http.Get(server2.URL + "/admin")
	fmt.Println("GET /admin -> status:", resp.StatusCode)
	resp.Body.Close()
}
