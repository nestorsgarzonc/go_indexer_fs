package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type EmailsZinc struct {
	Index  string  `json:"index"`
	Emails []Email `json:"records"`
}

type Email struct {
	From    string
	To      string
	Subject string
	Content string
}

const ZINC_URL = "http://localhost:4080"
const EMAIL_BASE_PATH = "./enron_mail_20110402/maildir"

func main() {
	fmt.Println("Welcome")
	emails, err := loadFolders()
	if err != nil {
		log.Fatal(err)
	}
	uploadToZinc(emails)
}

// Upload emails to Zinc in batches
func uploadToZinc(emails []Email) {
	jsonEmails, err := json.Marshal(EmailsZinc{
		Index:  "Subject",
		Emails: emails,
	})
	if err != nil {
		log.Fatal(err)
	}
	posturl := ZINC_URL + "/api/_bulkv2"
	// Create a HTTP post request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(jsonEmails))
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth("admin", "Complexpass#123")
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to upload emails to Zinc")
		fmt.Println("Response Status:", resp.Status)
		fmt.Println("Response Headers:", resp.Header)
		fmt.Println("Response Body:", resp.Body)
		log.Fatal("Failed to upload emails to Zinc")
	}
	fmt.Println("Emails uploaded to Zinc")
}

func loadFolders() ([]Email, error) {
	folders, err := os.ReadDir(EMAIL_BASE_PATH)
	if err != nil {
		return nil, err // Return the error instead of exiting
	}
	var wg sync.WaitGroup
	emailChan := make(chan []Email) // Create a channel to receive slices of emails
	for _, folder := range folders {
		if folder.IsDir() {
			wg.Add(1)
			go func(folderName string) {
				defer wg.Done()
				emails, err := loadFolder(folderName) // Make sure loadFolder returns ([]Email, error)
				if err != nil {
					log.Println(err) // Log the error; don't exit
					emailChan <- nil // Send nil to the channel to keep the flow
					return
				}
				emailChan <- emails // Send the result to the channel
			}(folder.Name())
		}
	}
	go func() {
		wg.Wait()
		close(emailChan) // Close the channel once all goroutines are done
	}()
	var allEmails []Email
	for emails := range emailChan {
		if emails != nil {
			allEmails = append(allEmails, emails...) // Collect all emails
		}
	}
	return allEmails, nil // Return the combined list of all emails
}

func loadFolder(path string) ([]Email, error) {
	var emails []Email
	subFolders, err := os.ReadDir(filepath.Join(EMAIL_BASE_PATH, path))
	if err != nil {
		return nil, err
	}
	for _, f := range subFolders {
		if f.IsDir() {
			files, err := os.ReadDir(filepath.Join(EMAIL_BASE_PATH, path, f.Name()))
			if err != nil {
				return nil, err
			}
			for _, fileX := range files {
				if !fileX.IsDir() {
					emailPath := filepath.Join(EMAIL_BASE_PATH, path, f.Name(), fileX.Name())
					email, err := loadEmail(emailPath)
					if err != nil {
						return nil, err
					}
					emails = append(emails, email)
				}
			}
		}
	}
	return emails, nil
}

func loadEmail(path string) (Email, error) {
	file, err := os.Open(path)
	if err != nil {
		return Email{}, err // Return the error instead of exiting
	}
	defer file.Close()
	var sb strings.Builder // Use strings.Builder for efficient concatenation
	email := Email{}       // Initialize with zero values
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	startContent := false
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "From:"):
			email.From = strings.TrimSpace(line[5:])
		case strings.HasPrefix(line, "To:"):
			email.To = strings.TrimSpace(line[3:])
		case strings.HasPrefix(line, "Subject:"):
			email.Subject = strings.TrimSpace(line[8:])
		case strings.HasPrefix(line, "X-FileName:"):
			startContent = true
		default:
			if startContent && len(line) > 0 {
				sb.WriteString(line + "\n")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Email{}, err // Handle potential errors from the scanner
	}
	email.Content = sb.String() // Assign the content from the builder
	return email, nil           // Return the email and a nil error
}
