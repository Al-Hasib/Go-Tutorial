package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

func main() {

	//a counter so /flaky can fail on purpose a couple of times, to
	//demonstrate retrying below
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/greet":
			name := r.URL.Query().Get("name")
			if name == "" {
				name = "stranger"
			}
			fmt.Fprintf(w, "hello, %s", name)
		case "/echo-header":
			fmt.Fprint(w, r.Header.Get("X-Token"))
		case "/slow":
			time.Sleep(300 * time.Millisecond)
			fmt.Fprint(w, "eventually done")
		case "/flaky":
			attempts++
			if attempts < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "success on attempt ", attempts)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// ---------- http.Get: THE SHORTCUT FOR SIMPLE GETs ----------

	resp, err := http.Get(server.URL + "/greet?name=Hasib")
	must(err)
	fmt.Println("GET /greet:", readAndClose(resp))

	// ---------- BUILDING A URL WITH QUERY PARAMETERS PROPERLY ----------

	//string-concatenating query params is fragile (spaces, &, = all need
	//escaping) - url.Values / Encode() handles that correctly
	q := url.Values{}
	q.Set("name", "A & B") // deliberately has characters that need escaping
	resp, err = http.Get(server.URL + "/greet?" + q.Encode())
	must(err)
	fmt.Println("GET with escaped query:", readAndClose(resp))

	// ---------- CUSTOM HEADERS NEED http.NewRequest + client.Do ----------

	req, err := http.NewRequest(http.MethodGet, server.URL+"/echo-header", nil)
	must(err)
	req.Header.Set("X-Token", "secret-123")
	resp, err = http.DefaultClient.Do(req)
	must(err)
	fmt.Println("GET with custom header:", readAndClose(resp))

	// ---------- TIMEOUTS: DON'T WAIT FOREVER ----------

	//a context deadline cancels the request if the server takes too long -
	//essential for any client talking to a network you don't fully control
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/slow", nil)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("GET /slow failed as expected:", err)
	}

	//http.Client{Timeout: ...} is the simpler, whole-request alternative -
	//it covers connect+read+everything, where a context lets you cancel for
	//other reasons too (e.g. the user navigated away)
	quickClient := &http.Client{Timeout: 50 * time.Millisecond}
	_, err = quickClient.Get(server.URL + "/slow")
	if err != nil {
		fmt.Println("client.Timeout also caught it:", err)
	}

	// ---------- CHECKING STATUS CODES, NOT JUST ERRORS ----------

	//err from Do()/Get() is about the REQUEST failing (network, timeout,
	//DNS...) - a 404 or 500 is still a "successful" HTTP round trip as far
	//as err is concerned. always check resp.StatusCode too.
	resp, err = http.Get(server.URL + "/does-not-exist")
	must(err)
	if resp.StatusCode >= 400 {
		fmt.Println("request succeeded but server returned an error status:", resp.StatusCode)
	}
	resp.Body.Close()

	// ---------- RETRYING A FLAKY ENDPOINT ----------

	for i := 1; i <= 5; i++ {
		resp, err = http.Get(server.URL + "/flaky")
		must(err)
		if resp.StatusCode == http.StatusOK {
			fmt.Println("flaky endpoint succeeded on try", i, "-", readAndClose(resp))
			break
		}
		resp.Body.Close()
		fmt.Println("try", i, "failed with status", resp.StatusCode, "- retrying")
		time.Sleep(time.Duration(i) * 10 * time.Millisecond) // simple growing backoff
	}
}

func readAndClose(resp *http.Response) string {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func must(err error) {
	if err != nil {
		panic(err) // lesson-only shortcut - real code should handle each error contextually
	}
}
