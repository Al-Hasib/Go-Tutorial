package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userStore struct {
	mu     sync.Mutex
	users  map[int]User
	nextID int
}

func newUserStore() *userStore {
	return &userStore{users: make(map[int]User), nextID: 1}
}

func (s *userStore) list() []User {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}

func (s *userStore) get(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	return u, ok
}

func (s *userStore) create(u User) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = s.nextID
	s.nextID++
	s.users[u.ID] = u
	return u
}

func (s *userStore) update(id int, u User) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return false
	}
	u.ID = id
	s.users[id] = u
	return true
}

func (s *userStore) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}

func main() {

	// ---------- REST: RESOURCES + VERBS + STATUS CODES ----------

	//REST maps CRUD operations onto HTTP methods against a URL that names a
	//RESOURCE (a collection "/users", or one item "/users/{id}") - the verb
	//says what to DO, the URL says what to do it TO.
	store := newUserStore()
	mux := http.NewServeMux()

	//Go 1.22+ ServeMux understands method + path pattern + {wildcard}
	//directly - no third-party router needed for this (see the dedicated
	//Routing lesson for more on patterns)
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, store.list())
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
			return
		}
		created := store.create(u)
		//201 Created (not 200) signals "a new resource now exists" - the
		//Location header points at exactly where to find it
		w.Header().Set("Location", fmt.Sprintf("/users/%d", created.ID))
		writeJSON(w, http.StatusCreated, created)
	})

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a number"})
			return
		}
		u, ok := store.get(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		writeJSON(w, http.StatusOK, u)
	})

	mux.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a number"})
			return
		}
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
			return
		}
		if !store.update(id, u) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		u.ID = id
		writeJSON(w, http.StatusOK, u) // PUT replaces the whole resource, so we return it
	})

	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a number"})
			return
		}
		if !store.delete(id) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		//204 No Content: successfully deleted, deliberately no body to return
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// ---------- WALKING THROUGH A TYPICAL CRUD LIFECYCLE ----------

	post(server.URL+"/users", `{"name":"Alice","email":"alice@example.com"}`)
	post(server.URL+"/users", `{"name":"Bob","email":"bob@example.com"}`)

	get(server.URL + "/users") // list: both users

	get(server.URL + "/users/1") // read one

	put(server.URL+"/users/1", `{"name":"Alice Smith","email":"alice.smith@example.com"}`) // update

	del(server.URL + "/users/1") // delete

	get(server.URL + "/users/1") // now 404 - gone

	del(server.URL + "/users/1") // deleting again is still a clean 404, not a crash -
	// this is what makes DELETE "idempotent": repeating it doesn't change the outcome
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func get(url string) {
	resp, _ := http.Get(url)
	report("GET "+url, resp)
}

func post(url, body string) {
	resp, _ := http.Post(url, "application/json", strings.NewReader(body))
	report("POST "+url, resp)
}

func put(url, body string) {
	req, _ := http.NewRequest(http.MethodPut, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	report("PUT "+url, resp)
}

func del(url string) {
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, _ := http.DefaultClient.Do(req)
	report("DELETE "+url, resp)
}

func report(label string, resp *http.Response) {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("%s -> %s %s\n", label, resp.Status, body)
}
