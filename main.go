package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	checkAPIKey()

	port := getPort()
	mux := http.NewServeMux()
	s := http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  readR,
		WriteTimeout: writeR,
		IdleTimeout:  keepA,
	}

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)

	log.Printf("Starting http server on port %s\n", port)
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server %s", err.Error())
		}
	}()
	log.Print("Server Started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Print("Signal closing server received")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Println("server shutdown failed", "error", err)
	}
	log.Println("server shutdown gracefully")
}
