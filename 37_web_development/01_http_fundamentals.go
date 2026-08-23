package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func main() {

	// ---------- HTTP: THE BASIC MODEL ----------

	//HTTP is a request/response protocol: a client sends one request, the
	//server sends back exactly one response, and that's the whole exchange -
	//nothing stays "in progress" afterward (more on this below, see STATELESSNESS).
	//
	//httptest.NewServer starts a real, listening HTTP server on a free local
	//port and hands back its URL - perfect for a self-contained lesson, no
	//manual port picking needed.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("server saw:", r.Method, r.URL.Path, "query:", r.URL.RawQuery)
		fmt.Println("  header User-Agent:", r.Header.Get("User-Agent"))
		fmt.Println("  header X-Demo:", r.Header.Get("X-Demo"))

		body, _ := io.ReadAll(r.Body)
		if len(body) > 0 {
			fmt.Println("  body:", string(body))
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK) // must be called before Write, and at most once
		w.Write([]byte("hello from the server"))
	}))
	defer server.Close()

	// ---------- METHODS CARRY MEANING ----------

	//GET: "give me this resource" - no body, safe to repeat, cacheable
	resp := mustDo("GET", server.URL+"/greeting?name=world", nil, nil)
	printResponse("GET", resp)

	//POST: "here's some data, do something with it" - has a body, and is
	//NOT safe to blindly repeat (e.g. don't auto-retry a POST that could
	//double-charge a payment)
	resp = mustDo("POST", server.URL+"/greeting", strings.NewReader("some data"), nil)
	printResponse("POST", resp)

	//PUT/DELETE/etc need a header or method net/http's shortcuts don't cover
	resp = mustDo("PUT", server.URL+"/greeting", nil, map[string]string{"X-Demo": "custom-header-value"})
	printResponse("PUT", resp)

	// ---------- STATUS CODES ARE GROUPED BY MEANING ----------

	fmt.Println()
	fmt.Println(http.StatusOK, http.StatusText(http.StatusOK))                                     // 2xx: success
	fmt.Println(http.StatusMovedPermanently, http.StatusText(http.StatusMovedPermanently))         // 3xx: redirect
	fmt.Println(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))                     // 4xx: client's fault
	fmt.Println(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))   // 5xx: server's fault

	// ---------- STATELESSNESS ----------

	//the handler above never stored anything about who called it - sending
	//the exact same request again gets treated identically, since the
	//server has no memory of the first one. recognizing "the same client"
	//across requests needs something extra (cookies, sessions, tokens) -
	//that's a later lesson (Authentication Concepts). HTTP itself gives you
	//none of that for free.
	fmt.Println()
	resp = mustDo("GET", server.URL+"/greeting?name=world", nil, nil)
	printResponse("GET again (server remembers nothing about the earlier requests)", resp)
}

func mustDo(method, url string, body io.Reader, headers map[string]string) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err) // lesson-only shortcut - real code should handle each error contextually
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func printResponse(label string, resp *http.Response) {
	defer resp.Body.Close() // always close a response body once you're done with it
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("%s -> status: %s, body: %q\n", label, resp.Status, string(body))
}
