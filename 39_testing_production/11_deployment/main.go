package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

// a minimal but real service - just enough to have something worth
// containerizing for this lesson: a normal endpoint, and a health check
// (see the Dockerfile's HEALTHCHECK, and the note on why this exists at all).
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // same env-var-with-default pattern as the Configuration lesson
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello from inside a container")
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		//a dedicated health endpoint is what an orchestrator (Docker,
		//Kubernetes, a load balancer) polls to decide "is this instance
		//alive and ready to receive traffic?" - keep it fast and dependency-light
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	slog.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
