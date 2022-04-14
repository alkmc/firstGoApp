package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
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

	log.Printf("starting http server on port: %s", port)
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()
	log.Printf("server started with %s", runtime.Version())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Print("signal closing server received")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}
	log.Print("server shutdown gracefully")
}
