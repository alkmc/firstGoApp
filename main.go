package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	tpl    = template.Must(template.ParseFiles("index.html"))
	apiKey *string
)

func main() {
	apiKey = flag.String("apikey", "", "Newsapi.org access key")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("apikey must be set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
