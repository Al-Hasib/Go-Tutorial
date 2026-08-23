package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type CreateOrderRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	Email     string `json:"email"`
}

// validate collects EVERY problem instead of stopping at the first one -
// far more useful to a client than fixing one field, resubmitting, and only
// then discovering a second problem.
func (r CreateOrderRequest) validate() []string {
	var errs []string
	if strings.TrimSpace(r.ProductID) == "" {
		errs = append(errs, "productId is required")
	}
	if r.Quantity <= 0 {
		errs = append(errs, "quantity must be greater than 0")
	}
	if r.Quantity > 100 {
		errs = append(errs, "quantity cannot exceed 100")
	}
	if !strings.Contains(r.Email, "@") {
		errs = append(errs, "email must be a valid email address")
	}
	return errs
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"errors": []string{"invalid JSON body"}})
			return
		}
		if errs := req.validate(); len(errs) > 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"errors": errs})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"status": "order placed", "order": req})
	})

	//validating a QUERY PARAMETER: it always arrives as a string, so
	//conversion + range checking both belong in validation
	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"errors": []string{"page must be a positive integer"}})
			return
		}
		fmt.Fprintf(w, "showing search results, page %d", page)
	})

	//validating a PATH PARAMETER the same way
	mux.HandleFunc("GET /orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"errors": []string{"id must be a number"}})
			return
		}
		fmt.Fprintf(w, "order #%d", id)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	post(server.URL+"/orders", `{"productId":"sku-1","quantity":2,"email":"a@b.com"}`)   // valid
	post(server.URL+"/orders", `{"productId":"","quantity":0,"email":"not-an-email"}`)   // 3 problems at once
	post(server.URL+"/orders", `{"productId":"sku-1","quantity":500,"email":"a@b.com"}`) // over the limit

	get(server.URL + "/search?page=2")
	get(server.URL + "/search?page=abc")
	get(server.URL + "/search?page=-1")

	get(server.URL + "/orders/42")
	get(server.URL + "/orders/not-a-number")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func post(url, body string) {
	resp, _ := http.Post(url, "application/json", strings.NewReader(body))
	report("POST "+url, resp)
}

func get(url string) {
	resp, _ := http.Get(url)
	report("GET "+url, resp)
}

func report(label string, resp *http.Response) {
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(label, "->", resp.Status, string(respBody))
}
