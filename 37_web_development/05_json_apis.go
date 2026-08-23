package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // "-" means "never serialize this field", period
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {

		// ---------- CHECKING Content-Type BEFORE TRUSTING THE BODY ----------

		//a client SAYING it's sending JSON (via this header) is different
		//from the body actually BEING valid JSON - checking this catches an
		//obviously wrong request before even attempting to parse anything
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
			writeErr(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		decoder := json.NewDecoder(r.Body)
		//DisallowUnknownFields rejects a request containing fields the
		//struct doesn't have - catches typos ("usernme") and unexpected
		//extra data instead of silently ignoring them
		decoder.DisallowUnknownFields()

		var req SignupRequest
		if err := decoder.Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, describeJSONError(err))
			return
		}
		if req.Username == "" || req.Password == "" {
			writeErr(w, http.StatusBadRequest, "username and password are both required")
			return
		}

		//PasswordHash is never filled from user input, and is tagged
		//json:"-" so it can NEVER accidentally leak back out in a response,
		//no matter what any handler later does with this struct
		resp := UserResponse{Username: req.Username, PasswordHash: "(hashed elsewhere)"}
		writeJSON(w, http.StatusCreated, resp)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// ---------- A VALID REQUEST ----------
	post(server.URL+"/signup", "application/json", `{"username":"hasib","password":"secret"}`)

	// ---------- WRONG Content-Type ----------
	post(server.URL+"/signup", "text/plain", `{"username":"hasib","password":"secret"}`)

	// ---------- MALFORMED JSON (SYNTAX ERROR) ----------
	post(server.URL+"/signup", "application/json", `{"username":"hasib", "password": secret}`) // unquoted value - invalid JSON

	// ---------- WRONG TYPE FOR A FIELD ----------
	post(server.URL+"/signup", "application/json", `{"username":"hasib","password":12345}`)

	// ---------- UNKNOWN FIELD ----------
	post(server.URL+"/signup", "application/json", `{"username":"hasib","password":"secret","isAdmin":true}`)

	// ---------- MISSING FIELD ----------
	post(server.URL+"/signup", "application/json", `{"username":"hasib"}`)
}

// describeJSONError turns encoding/json's error types into a clearer
// message than the raw error text - real APIs often map errors like this
func describeJSONError(err error) string {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxErr):
		return "request body is not valid JSON"
	case errors.As(err, &typeErr):
		return fmt.Sprintf("field %q has the wrong type", typeErr.Field)
	case strings.Contains(err.Error(), "unknown field"):
		return err.Error() // already descriptive: `json: unknown field "isAdmin"`
	default:
		return "could not parse request body"
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func post(url, contentType, body string) {
	resp, _ := http.Post(url, contentType, strings.NewReader(body))
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("POST body=%s\n  -> %s %s\n", body, resp.Status, respBody)
}
