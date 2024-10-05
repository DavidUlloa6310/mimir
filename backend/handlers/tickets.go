package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/davidulloa/mimir/models"
)

type TicketRequestBody struct {
    InstanceID string `json:"instance_id"`
}

type TicketResponseBody struct {
    Result []models.Ticket `json:"result"`
}

type TicketHandler struct {
    Client *http.Client
}

// NewTicketHandler creates a new instance of the TicketHandler
func NewTicketHandler(client *http.Client) *TicketHandler {
    return &TicketHandler{
        Client: client,
    }
}

func (h *TicketHandler) TicketsHandler(w http.ResponseWriter, r *http.Request) {
    username, password, ok := r.BasicAuth()
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Unable to read request body", http.StatusBadRequest)
        return
    }

    var responseBody TicketRequestBody
    err = json.Unmarshal(body, &responseBody)
    if err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    apiURL := fmt.Sprintf("https://%s.service-now.com/api/now/table/incident", responseBody.InstanceID)

    queryParams := url.Values{}
    queryParams.Add("sysparm_limit", "10")
    queryParams.Add("sysparm_fields", "number,short_description,priority,state")

    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        log.Fatalf("Error creating request: %v", err)
    }

    req.URL.RawQuery = queryParams.Encode()
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.SetBasicAuth(username, password)

    resp, err := h.Client.Do(req)
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

func jsonResponse(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
