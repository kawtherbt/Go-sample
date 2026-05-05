package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Contact struct
type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func triggerPipeline(c Contact) error {
	url := "https://api.github.com/repos/kawtherbt/Go-sample/dispatches"

	// Payload with form data
	payloadMap := map[string]interface{}{
		"event_type": "form_submitted",
		"client_payload": map[string]string{
			"name":    c.Name,
			"email":   c.Email,
			"message": c.Message,
		},
	}

	payloadBytes, _ := json.Marshal(payloadMap)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN not set")
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	return nil
}

// form sub
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
		c.Name, c.Email, c.Message) //log

	err = triggerPipeline(c)
	if err != nil {
		log.Println("Error triggering pipeline:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message received, but failed to trigger pipeline"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Message received and pipeline triggered 🚀"))
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/contact", contactHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}