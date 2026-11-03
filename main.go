package main

import (
	"errors"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve files from the root directory as if it was /app/
	mux.Handle(
		"/app/",
		http.StripPrefix(
			"/app",
			http.FileServer(http.Dir(".")),
		),
	)

	// Readiness endpoint, used by a client to check if this server is ready
	// to process requests.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

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
