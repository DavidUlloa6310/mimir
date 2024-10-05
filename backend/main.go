package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Ticket represents a ticket structure
type Ticket struct {
	Number           int    `json:"number"`
	ShortDescription string `json:"short_description"`
	State            string `json:"state"`
	Priority         string `json:"priority"`
}

type Accelerator struct {
	ID    int    `json:"id"`
	Url   string `json:"url"`
	Title string `json:"title"`
	// Service string `json:"string"`
}


// Suggestion represents a suggestion structure
type Suggestion struct {
	ID             int            `json:"id"`
	Description    string         `json:"description"`
	Title          string         `json:"title"`
	Accelerators   []Accelerator  `json:"accelerators"`
}

// Documentation represents a documentation entry
type Documentation struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ChatMessage struct {
	Message string `json:"message"`
	Role string `json:"role"`
}

type TicketRequestBody struct {
	InstanceID string `json:"instance_id"`
}

type TicketResponseBody struct {
	Result []Ticket `json:"result"`
}

func ticketsHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// Parse the JSON body
	var responseBody TicketRequestBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	apiURL := fmt.Sprintf("%s/api/now/table/incident", responseBody.InstanceID)

	// Define any query parameters you need (optional)
	queryParams := url.Values{}
	queryParams.Add("sysparm_limit", "10")
	queryParams.Add("sysparm_fields", "number,short_description,priority,state")

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.URL.RawQuery = queryParams.Encode()

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	tickets := &TicketRequestBody{}
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		err = json.Unmarshal(body, tickets)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}

		fmt.Printf("Incident Data: %v\n", tickets)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to retrieve data. Status code: %d, Response: %s\n", resp.StatusCode, string(body))
	}

	jsonResponse(w, tickets)
}

// suggestionsHandler returns dummy suggestion data
func suggestionsHandler(w http.ResponseWriter, r *http.Request) {
	suggestions := []Suggestion{}

	jsonResponse(w, suggestions)
}

// chatHandler returns dummy chat messages
func chatHandler(w http.ResponseWriter, r *http.Request) {
	chat := []ChatMessage{
		{Role: "John", Message: "Hello, how can I help?"},
		{Role: "Jane", Message: "I need assistance with my order."},
	}

	jsonResponse(w, chat)
}

// documentationHandler returns dummy documentation entries
func documentationHandler(w http.ResponseWriter, r *http.Request) {
	documentation := []Documentation{
		{Title: "API Guide", Content: "This is a guide for using our API."},
		{Title: "Getting Started", Content: "How to get started with our platform."},
	}

	jsonResponse(w, documentation)
}

// jsonResponse is a helper function to send JSON responses
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Register handlers for each route
	http.HandleFunc("/tickets", ticketsHandler)
	http.HandleFunc("/suggestions", suggestionsHandler)
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/documentation", documentationHandler)

	// Start the server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
