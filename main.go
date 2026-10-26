package main

import (
	"errors"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve current directory, generally a bad idea
	mux.Handle("/", http.FileServer(http.Dir(".")))

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
