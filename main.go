package main

import (
	"fmt"
	"log"
	"net/http"

	"test-service/config"
	"test-service/handlers"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("/one",   handlers.One)
	mux.HandleFunc("/two",   handlers.Two)
	mux.HandleFunc("/three", handlers.Three)
	mux.HandleFunc("/four",  handlers.Four)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
