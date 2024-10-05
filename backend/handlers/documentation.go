package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/davidulloa/mimir/models"
)

type DocumentationHandler struct {}

// NewDocumentationHandler creates and returns a new DocumentationHandler instance
func NewDocumentationHandler() *DocumentationHandler {
	return &DocumentationHandler{}
}

// DocumentationHandler serves the documentation route and returns dummy data
func (h *DocumentationHandler) DocumentationHandler(w http.ResponseWriter, r *http.Request) {
	// Dummy documentation data
	documentation := []models.Documentation{
		{Title: "API Guide", Content: "This is a guide for using our API."},
		{Title: "Getting Started", Content: "How to get started with our platform."},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documentation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}