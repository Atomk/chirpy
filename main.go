package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	// Counts requests involving the file server.
	fileserverHits atomic.Int32
}

func main() {
	cfg := apiConfig{}
	mux := http.NewServeMux()

	// Serve files from the root directory as if it was /app/
	mux.Handle(
		"/app/",
		http.StripPrefix(
			"/app",
			// Increments a counter each time the FileServer handler is used.
			cfg.middlewareMetricsInc(http.FileServer(http.Dir("."))),
		),
	)

	// Readiness endpoint, used by a client to check if this server is ready
	// to process requests.
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /metrics", cfg.getNumberOfRequests)
	mux.HandleFunc("POST /reset", cfg.resetMetrics)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Printf("launching the server on %s\n", server.Addr)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server was shut down")
	} else {
		fmt.Printf("the server was unexpectedly closed: %e\n", err)
	}
}

// This middleware is a sort of stateful handler: returns the given handled,
// wrapped into a function that when executed increments a counter.
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getNumberOfRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	text := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(text))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
}
