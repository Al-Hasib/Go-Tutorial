package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func main() {

	mux := http.NewServeMux()

	// ---------- STATIC PATHS ----------

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// ---------- METHOD + PATH TOGETHER ----------

	//the exact same path can have a different handler per method - built
	//into the pattern syntax since Go 1.22, no manual r.Method switch needed
	mux.HandleFunc("GET /articles", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "list of articles")
	})
	mux.HandleFunc("POST /articles", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "created an article")
	})

	// ---------- WILDCARDS ----------

	mux.HandleFunc("GET /articles/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "article #%s", r.PathValue("id"))
	})

	//multiple wildcards in one pattern
	mux.HandleFunc("GET /articles/{id}/comments/{commentID}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "comment #%s on article #%s", r.PathValue("commentID"), r.PathValue("id"))
	})

	// ---------- PRECEDENCE: A MORE SPECIFIC PATTERN WINS ----------

	//"/articles/latest" and "/articles/{id}" could both match the path
	//"/articles/latest" - ServeMux always prefers the more specific (fewer
	//wildcards) pattern, regardless of which one was REGISTERED first in
	//the source code - precedence is based on pattern specificity, not order
	mux.HandleFunc("GET /articles/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "the latest article (never falls through to {id})")
	})

	// ---------- {name...} MATCHES THE REST OF THE PATH ----------

	//useful for things like serving files at arbitrary nested paths
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "serving file at path: %s", r.PathValue("path"))
	})

	// ---------- SUBTREE (TRAILING SLASH) PATTERNS ----------

	//a pattern ending in "/" matches that path AND everything nested under it
	mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "static asset: %s", r.URL.Path)
	})

	// ---------- CUSTOM 404 ----------

	//ServeMux has no direct "SetNotFound" - the standard trick is a
	//catch-all pattern at "/", which only fires when nothing more specific matched
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "custom 404: nothing registered for %s", r.URL.Path)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	get(server.URL + "/health")
	get(server.URL + "/articles")
	get(server.URL + "/articles/42")
	get(server.URL + "/articles/42/comments/7")
	get(server.URL + "/articles/latest") // proves the specific pattern wins over {id}
	get(server.URL + "/files/images/logo.png")
	get(server.URL + "/static/app.css")
	get(server.URL + "/this-route-does-not-exist")
}

func get(url string) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("GET %s -> %s: %s\n", url, resp.Status, body)
}
