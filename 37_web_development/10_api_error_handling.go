package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

// APIError is a typed error carrying the HTTP status it should map to -
// letting handlers return ONE kind of error and have status codes stay
// consistent across the whole API, instead of every handler deciding
// ad-hoc which WriteHeader to call.
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string { return e.Message }

func NotFound(msg string) *APIError { return &APIError{Status: http.StatusNotFound, Message: msg} }

var errOutOfStock = errors.New("product is out of stock") // a plain sentinel error from deeper in the code

// ---------- THE ERROR RESPONSE SHAPE ----------

// every error response looks the same: {"error": {"message": "..."}} -
// clients only ever need to learn ONE shape, instead of a different one per endpoint.
type errorEnvelope struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	var env errorEnvelope
	env.Error.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(env)
}

// apiHandler lets handlers just RETURN an error instead of writing the
// response themselves - one place then decides how every kind of error
// becomes an HTTP response.
type apiHandler func(w http.ResponseWriter, r *http.Request) error

func wrap(h apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			//a panic anywhere in the handler becomes a clean 500 instead of
			//taking down request handling - same idea as the Middleware lesson's recovery
			if rec := recover(); rec != nil {
				log.Printf("panic in %s %s: %v", r.Method, r.URL.Path, rec)
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		err := h(w, r)
		if err == nil {
			return
		}

		var apiErr *APIError
		switch {
		case errors.As(err, &apiErr):
			//a KNOWN, expected error - safe to show its message to the
			//client, since whoever wrote it chose that message on purpose
			writeError(w, apiErr.Status, apiErr.Message)
		case errors.Is(err, errOutOfStock):
			writeError(w, http.StatusConflict, err.Error())
		default:
			//an UNEXPECTED error - log the real details server-side (for
			//debugging), but never send them to the client: a raw Go error
			//or stack trace can leak internal details (file paths, SQL,
			//internal type names) that are none of the client's business
			log.Printf("unexpected error in %s %s: %v", r.Method, r.URL.Path, err)
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products/{id}", wrap(func(w http.ResponseWriter, r *http.Request) error {
		id := r.PathValue("id")
		if id != "42" {
			return NotFound("product " + id + " does not exist") // 4xx: client asked for something that isn't there
		}
		fmt.Fprint(w, `{"id":"42","name":"widget"}`)
		return nil
	}))

	mux.HandleFunc("POST /orders", wrap(func(w http.ResponseWriter, r *http.Request) error {
		return errOutOfStock // simulates an error bubbling up from deeper business logic
	}))

	mux.HandleFunc("GET /crash", wrap(func(w http.ResponseWriter, r *http.Request) error {
		var m map[string]int
		m["this"] = 1 // panics: assignment to entry in nil map
		return nil
	}))

	mux.HandleFunc("GET /mystery", wrap(func(w http.ResponseWriter, r *http.Request) error {
		//5xx, and the real message below is deliberately NOT what reaches the client
		return fmt.Errorf("unclassified internal failure, details: db=prod-1 query=SELECT...")
	}))

	server := httptest.NewServer(mux)
	defer server.Close()

	get(server.URL + "/products/42")
	get(server.URL + "/products/99")
	get(server.URL + "/orders") // registered only for POST -> ServeMux auto-returns 405
	post(server.URL + "/orders")
	get(server.URL + "/crash")
	get(server.URL + "/mystery")
}

func get(url string) {
	resp, _ := http.Get(url)
	report("GET "+url, resp)
}

func post(url string) {
	resp, _ := http.Post(url, "application/json", nil)
	report("POST "+url, resp)
}

func report(label string, resp *http.Response) {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(label, "->", resp.Status, string(body))
}
