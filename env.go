package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"time"
)

var (
	tpl    = template.Must(template.ParseFiles("index.html"))
	apiKey *string
)

const (
	readR   = 5 * time.Second   // max time to read request from the client
	writeR  = 10 * time.Second  // max time to write response to the client
	keepA   = 120 * time.Second // max time for connections using TCP Keep-Alive
	timeout = 10 * time.Second  // max time to complete tasks before shutdown
	port    = ":3000"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return port
	}
	return ":" + port
}

func checkAPIKey() {
	apiKey = flag.String("apiKey", "", "newsapi.org access key")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("apiKey must be set")
	}
}
