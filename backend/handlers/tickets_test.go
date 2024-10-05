package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTicketsHandler(t *testing.T) {
    client := &http.Client{}
    handler := NewTicketHandler(client)

    // Prepare a request body
    requestBody, _ := json.Marshal(TicketRequestBody{InstanceID: "testInstance"})
    req, err := http.NewRequest("POST", "/tickets", bytes.NewBuffer(requestBody))
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }

    req.SetBasicAuth("testuser", "testpass")

    // Create a ResponseRecorder to capture the response
    rr := httptest.NewRecorder()

    handler.TicketsHandler(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}
