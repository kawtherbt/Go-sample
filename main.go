package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v85/github"
	"golang.org/x/oauth2"
)

// Contact struct (form data)
type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Trigger GitHub Actions pipeline using go-github
func triggerPipeline(c Contact) error {
	ctx := context.Background()

	// 1. Setup Authentication
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("TOKEN")}, // 👈 keep env, not hardcoded
	)
	tc := oauth2.NewClient(ctx, ts)

	// 2. Create the GitHub Client
	client := github.NewClient(tc)

	// 3. Trigger repository_dispatch (instead of Get)
	payload := map[string]interface{}{
	"name":    c.Name,
	"email":   c.Email,
	"message": c.Message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	dispatch := github.DispatchRequestOptions{
		EventType:     "form_submitted",
		ClientPayload: (*json.RawMessage)(&jsonPayload),
	}

	_,_, err = client.Repositories.Dispatch(ctx, "kawtherbt", "Go-sample", dispatch)
	if err != nil {
		return fmt.Errorf("error triggering pipeline: %v", err)
	}

	return nil
}

// Handle form submission
func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var c Contact
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log form submission
	fmt.Printf("New message:\nName: %s\nEmail: %s\nMessage: %s\n\n",
		c.Name, c.Email, c.Message)

	// Trigger pipeline
	err = triggerPipeline(c)
	if err != nil {
		log.Println("Error triggering pipeline:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message received, but failed to trigger pipeline"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Message received and pipeline triggered"))
}

func main() {
	// Serve static files (HTML, CSS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoint
	http.HandleFunc("/contact", contactHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}