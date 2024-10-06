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

type TicketResponseBody struct {
    Result []models.Ticket `json:"result"`
}

type TicketHandler struct {
    Client *http.Client
}

type IncidentsApiResponse struct {
	Result []models.Incident `json:"result"`
}

// NewTicketHandler creates a new instance of the TicketHandler
func NewTicketHandler(client *http.Client) *TicketHandler {
    return &TicketHandler{
        Client: client,
    }
}

func GetIncidents(client *http.Client, instanceID string, username string, password string) *IncidentsApiResponse {
    apiURL := fmt.Sprintf("https://%s.service-now.com/api/now/table/incident", instanceID)

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

    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error making request: %v", err)
    }
    defer resp.Body.Close()

    incidents := &IncidentsApiResponse{}
    if resp.StatusCode == http.StatusOK {
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Fatalf("Error reading response body: %v", err)
        }

        err = json.Unmarshal(body, incidents)
        if err != nil {
            log.Fatalf("Error parsing JSON: %v", err)
        }
    } else {
        body, _ := io.ReadAll(resp.Body)
        fmt.Printf("Failed to retrieve data. Status code: %d, Response: %s\n", resp.StatusCode, string(body))
    }
    return incidents
}

func ToTickets(incidents *IncidentsApiResponse) []models.Ticket {
	tickets := []models.Ticket{}
	for _, incident := range incidents.Result {
		tickets = append(tickets, models.Ticket{
			ShortDescription: incident.ShortDescription,	
			Priority: incident.Priority,
			Number: incident.Number,
			State: incident.State,
		})
	}
    return tickets
}

func (h *TicketHandler) TicketsHandler(w http.ResponseWriter, r *http.Request) {
    instanceID, username, password, err := ParseCredentials(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    incidents := GetIncidents(h.Client, instanceID, username, password)
    tickets := ToTickets(incidents)
    jsonResponse(w, tickets)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
