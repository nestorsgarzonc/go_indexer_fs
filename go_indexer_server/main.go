package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const ZINC_URL = "http://localhost:4080"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/emails", fetchEmails)
		r.Get("/hello", helloHandler)
	})
	fmt.Println("Starting server at port 8080")
	http.ListenAndServe(":8080", r)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func fetchEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	query := r.Form.Get("query")
	fmt.Println("Query: " + query)
	emails, err := zincQuery(query)
	if err != nil {
		http.Error(w, "Failed to fetch emails from Zinc", http.StatusInternalServerError)
		return
	}
	emailsJson, err := json.Marshal(emails)
	if err != nil {
		http.Error(w, "Failed to marshal emails", http.StatusInternalServerError)
		return
	}
	w.Write(emailsJson)
}

func zincQuery(query string) ([]Source, error) {
	getUrl := ZINC_URL + "/api/Subject/_search"
	var body string
	if query == "" {
		body = `{"sort_fields": ["-@timestamp"], "from": 0, "max_results": 20, "_source": []}`
	} else {
		body = `
		{
    "search_type": "match",
    "query": { "term": ` + query + `, "field": "_all" },
    "sort_fields": ["-@timestamp"],
    "from": 0,
    "max_results": 20,
    "_source": []}`
	}
	r, err := http.NewRequest("POST", getUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth("admin", "Complexpass#123")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Response Status:", resp.Status)
		fmt.Println("Response Headers:", resp.Header)
		fmt.Println("Response Body:", resp.Body)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	email, err := UnmarshalEmail(bodyBytes)
	if err != nil {
		return nil, err
	}
	emails := email.Hits.Hits
	filteredEmails := make([]Source, len(emails))
	for i, email := range emails {
		filteredEmails[i] = email.Source
	}
	return filteredEmails, nil
}
