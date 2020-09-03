package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	params := u.Query()
	searchKey := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	s := &searchNews{}
	s.SearchKey = searchKey

	next, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Unexpected server error", http.StatusInternalServerError)
		return
	}
	s.NextPage = next

	const (
		URL      = "https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%d&apiKey=%s&sortBy=publishedAt&language=en"
		pageSize = 20
	)

	endpoint := fmt.Sprintf(URL, url.QueryEscape(s.SearchKey), pageSize, s.NextPage, *apiKey)

	if err := fetch(endpoint, &s.Results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	s.TotalPages = totalPages(s.Results.TotalResults, pageSize)

	if ok := !s.IsLastPage(); ok {
		s.NextPage++
	}
	if err := tpl.Execute(w, s); err != nil {
		log.Fatal(err)
	}
}

func fetch(endpoint string, v interface{}) error {
	resp, err := http.Get(endpoint)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return errors.New("Could not fetch data")
	}
	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		newsErr := &newsAPIError{}
		err := dec.Decode(newsErr)
		if errs := decErr(err); errs != nil {
			return errs
		}
		return errors.New(newsErr.Message)
	}

	err = dec.Decode(v)
	if errs := decErr(err); errs != nil {
		return errs
	}

	return nil
}

func decErr(err error) error {
	if err != nil {
		log.Println(err.Error())
		return errors.New("JSON decoding error")
	}
	return nil
}

func totalPages(total, pageSize int) int {
	return int(math.Ceil(float64(total / pageSize)))
}
