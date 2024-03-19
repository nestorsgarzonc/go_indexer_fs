package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const ZINC_URL = "http://localhost:4080"

func main() {
	http.HandleFunc("/emails", fetchEmails)
	http.HandleFunc("/hello", helloHandler)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Hello!")
}

func fetchEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.URL.Path != "/emails" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
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
