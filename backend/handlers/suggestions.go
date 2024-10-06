package handlers

import (
	"encoding/json"
	"net/http"
)

// SuggestionsHandler handles suggestions-related requests
type SuggestionsHandler struct {
	// You can add dependencies here, such as database connections, services, etc.
}

type SuggestionsBody struct {
	ticketIds []string `json:"tickets"`
}

// NewSuggestionsHandler creates and returns a new SuggestionsHandler instance
func NewSuggestionsHandler() *SuggestionsHandler {
	return &SuggestionsHandler{}
}

// SuggestionsHandler serves the suggestions route and returns dummy data
func (h *SuggestionsHandler) SuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	// Dummy suggestion data
	var data SuggestionsBody
	err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	instanceId, username, password, err := ParseCredentials(r)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	incidents := GetIncidents(client, instanceId, username, password)

	tickets := ToTickets(incidents)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}