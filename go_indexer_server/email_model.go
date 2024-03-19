package main

import "encoding/json"

func UnmarshalEmail(data []byte) (Email, error) {
	var r Email
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Email) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Email struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Hits []Hit `json:"hits"`
}

type Hit struct {
	Source Source `json:"_source"`
}

type Source struct {
	Timestamp string `json:"@timestamp"`
	Content   string `json:"Content"`
	From      string `json:"From"`
	Subject   string `json:"Subject"`
	To        string `json:"To"`
}
