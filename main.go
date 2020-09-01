package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	tpl    = template.Must(template.ParseFiles("index.html"))
	apiKey *string
)

type source struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

type article struct {
	Source      source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

func (a *article) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}

type results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []article `json:"articles"`
}

type search struct {
	SearchKey  string
	NextPage   int
	TotalPages int
	Results    results
}

func (s *search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

func (s *search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}
	return s.NextPage - 1
}

func (s *search) PreviousPage() int {
	return s.CurrentPage() - 1
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

type newsAPIError struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	searchKey := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	search := &search{}
	search.SearchKey = searchKey

	next, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Unexpected server error", http.StatusInternalServerError)
		return
	}

	search.NextPage = next
	pageSize := 20

	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%d&apiKey=%s&sortBy=publishedAt&language=en",
		url.QueryEscape(search.SearchKey), pageSize, search.NextPage, *apiKey)
	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		newErr := &newsAPIError{}
		err := json.NewDecoder(resp.Body).Decode(newErr)
		if err != nil {
			http.Error(w, "Unexpected server error", http.StatusInternalServerError)
			return
		}
		http.Error(w, newErr.Message, http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&search.Results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	search.TotalPages = int(math.Ceil(float64(search.Results.TotalResults / pageSize)))

	if ok := !search.IsLastPage(); ok {
		search.NextPage++
	}

	err = tpl.Execute(w, search)
	if err != nil {
		log.Fatal(err)
	}

}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, runtime.Version())
}

func main() {
	apiKey = flag.String("apikey", "", "Newsapi.org access key")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("apiKey must be set")
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
	mux.HandleFunc("/version", versionHandler)
	http.ListenAndServe(":"+port, mux)
}
