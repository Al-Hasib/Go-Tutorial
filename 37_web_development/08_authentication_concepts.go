package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

// ---------- STORING PASSWORDS: THE RIGHT IDEA, THE WRONG TOOL ----------

// hashing a password before storing it (never the raw password) is the
// right IDEA. sha256 is only used here to keep the mechanics simple - real
// password storage should use bcrypt or argon2 (e.g. golang.org/x/crypto/bcrypt),
// which are deliberately SLOW and add a random salt per password. plain
// sha256 is fast, which is exactly what makes it a bad fit for passwords:
// an attacker with a leaked hash database can try billions of guesses per
// second against it.
func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

var userPasswordHashes = map[string]string{
	"hasib": hashPassword("correct-password"),
}

var apiKeys = map[string]string{
	"sk_live_abc123": "service-account-1",
}

// standing in for what would really be a database or cache - maps an
// issued token to the username it belongs to
var activeTokens = map[string]string{}

func main() {

	mux := http.NewServeMux()

	// ---------- HTTP BASIC AUTH ----------

	//credentials travel in the Authorization header, base64-encoded (NOT
	//encrypted - Basic Auth is only safe over HTTPS) as "username:password"
	mux.HandleFunc("GET /basic", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		expectedHash, exists := userPasswordHashes[username]

		//subtle.ConstantTimeCompare avoids a timing attack: a plain "=="
		//comparison can return faster on an early mismatched byte, which an
		//attacker could measure to guess the hash one byte at a time
		if !ok || !exists || subtle.ConstantTimeCompare([]byte(hashPassword(password)), []byte(expectedHash)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			w.WriteHeader(http.StatusUnauthorized) // 401: "who are you?" - no/bad credentials
			return
		}
		fmt.Fprintf(w, "welcome, %s", username)
	})

	// ---------- BEARER TOKENS ----------

	//a token issued once (e.g. after a successful login) is sent on every
	//later request instead of a password - the standard shape for APIs,
	//since the client never has to keep resending the raw password
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		expectedHash, exists := userPasswordHashes[username]
		if !exists || subtle.ConstantTimeCompare([]byte(hashPassword(password)), []byte(expectedHash)) != 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := "tok_" + username + "_12345" // a real system uses a long random value here
		activeTokens[token] = username
		fmt.Fprint(w, token)
	})

	mux.HandleFunc("GET /profile", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(auth) <= len(prefix) || auth[:len(prefix)] != prefix {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		username, ok := activeTokens[auth[len(prefix):]]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "profile for %s", username)
	})

	// ---------- API KEYS ----------

	//simpler than user login - meant for service-to-service or third-party
	//integrations, not for representing a logged-in person
	mux.HandleFunc("GET /service-data", func(w http.ResponseWriter, r *http.Request) {
		owner, ok := apiKeys[r.Header.Get("X-API-Key")]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "service data for %s", owner)
	})

	// ---------- 401 VS 403 ----------

	//401 Unauthorized: "I don't know who you are" (missing/invalid credentials)
	//403 Forbidden: "I know exactly who you are, and you're not allowed to do this"
	mux.HandleFunc("GET /admin-only", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(auth) <= len(prefix) {
			w.WriteHeader(http.StatusUnauthorized) // no credentials at all
			return
		}
		username, ok := activeTokens[auth[len(prefix):]]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized) // credentials don't correspond to anyone
			return
		}
		if username != "admin" {
			w.WriteHeader(http.StatusForbidden) // known identity, just not allowed here
			return
		}
		fmt.Fprint(w, "admin dashboard")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// ---------- BASIC AUTH DEMO ----------

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/basic", nil)
	req.SetBasicAuth("hasib", "correct-password")
	resp, _ := http.DefaultClient.Do(req)
	report("GET /basic (correct creds)", resp)

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/basic", nil)
	req.SetBasicAuth("hasib", "wrong-password")
	resp, _ = http.DefaultClient.Do(req)
	report("GET /basic (wrong password)", resp)

	// ---------- BEARER TOKEN DEMO ----------

	req, _ = http.NewRequest(http.MethodPost, server.URL+"/login", nil)
	req.SetBasicAuth("hasib", "correct-password")
	resp, _ = http.DefaultClient.Do(req)
	token := readBody(resp)
	fmt.Println("issued token:", token)

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = http.DefaultClient.Do(req)
	report("GET /profile (with token)", resp)

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/profile", nil)
	req.Header.Set("Authorization", "Bearer garbage-token")
	resp, _ = http.DefaultClient.Do(req)
	report("GET /profile (bad token)", resp)

	// ---------- API KEY DEMO ----------

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/service-data", nil)
	req.Header.Set("X-API-Key", "sk_live_abc123")
	resp, _ = http.DefaultClient.Do(req)
	report("GET /service-data (valid key)", resp)

	// ---------- 401 VS 403 DEMO ----------

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/admin-only", nil)
	resp, _ = http.DefaultClient.Do(req)
	report("GET /admin-only (no token)", resp) // 401

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+token) // valid token, but user is "hasib", not "admin"
	resp, _ = http.DefaultClient.Do(req)
	report("GET /admin-only (logged in, not admin)", resp) // 403
}

func report(label string, resp *http.Response) {
	defer resp.Body.Close()
	fmt.Println(label, "->", resp.Status)
}

func readBody(resp *http.Response) string {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
