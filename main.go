package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"bytes"
)

type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var c Contact
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Printf("New message:\nName: %s\nEmail: %s\nMessage: %s\n\n",
		c.Name, c.Email, c.Message)

	w.Write([]byte("Message received  "))
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/contact", contactHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func triggerPipeline() {
	url := "https://api.github.com/repos/kawtherbt/Go-sample/dispatches"

	payload := []byte(`{
		"event_type": "form_submitted"
	}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer YOUR_GITHUB_TOKEN")

	client := &http.Client{}
	client.Do(req)
}