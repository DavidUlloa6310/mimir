package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidulloa/mimir/models"
)

func TestTicketsHandler(t *testing.T) {
	// Sample JSON response to be returned by the mock HTTP client
	testResponse := `{
		"result": [
			{
				"number": "INC0010110",
				"short_description": "Jira sprint planning feature is glitchy, we lost all story points for a session.",
				"priority": "5",
				"state": "1"
			}
		]
	}`

	client := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {
			// Optionally, check the request details here (e.g., URL, headers)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(testResponse)),
				Header:     make(http.Header),
			}
		}),
	}

	handler := NewTicketHandler(client)

	requestBody := `{"instance_id": "test_instance"}`
	req := httptest.NewRequest("POST", "/tickets", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("testuser", "testpass")

	rr := httptest.NewRecorder()

	handler.TicketsHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expectedTickets := []models.Ticket{
		{
			Number:           "INC0010110",
			ShortDescription: "Jira sprint planning feature is glitchy, we lost all story points for a session.",
			Priority:         "5",
			State:            "1",
		},
	}

	var actualTickets []models.Ticket
	err := json.Unmarshal(rr.Body.Bytes(), &actualTickets)
	if err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}

	if len(actualTickets) != len(expectedTickets) {
		t.Fatalf("Expected %d tickets, got %d", len(expectedTickets), len(actualTickets))
	}

	for i, ticket := range actualTickets {
		if ticket != expectedTickets[i] {
			t.Errorf("Ticket %d does not match expected result.\nGot: %+v\nWant: %+v", i, ticket, expectedTickets[i])
		}
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
