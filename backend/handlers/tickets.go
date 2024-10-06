package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/davidulloa/mimir/database"
	"github.com/davidulloa/mimir/models"
)

type TicketResponseBody struct {
    Result []models.Ticket `json:"result"`
}

type TicketCache struct {
    mu             sync.RWMutex
    lastIncidents  IncidentsApiResponse
    lastClusters   database.TicketResponse
    lastUpdateTime time.Time
}

type TicketHandler struct {
    Client *http.Client
    Cache  *TicketCache
}

type IncidentsApiResponse struct {
	Result []models.Incident `json:"result"`
}

type ClusteredTickets struct {
    ClusterDescription string   `json:"cluster_description"`
    Tickets            []models.Ticket `json:"tickets"`
}

type ClusteredTicketResponse struct {
    Clusters []ClusteredTickets `json:"clusters"`
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
        return
    }

    incidents := GetIncidents(h.Client, instanceID, username, password)

    h.Cache.mu.RLock()
    cacheValid := reflect.DeepEqual(incidents, h.Cache.lastIncidents) && time.Since(h.Cache.lastUpdateTime) < 5*time.Minute
    h.Cache.mu.RUnlock()

    var clusters database.TicketResponse

    // Caching TF-IDF Clusters
    if cacheValid {
        h.Cache.mu.RLock()
        clusters = h.Cache.lastClusters
        h.Cache.mu.RUnlock()
    } else {
        tickets := ToTickets(incidents)
        database.StoreTickets(tickets)
        shortDescriptions := make([]string, len(tickets))
        for i, ticket := range tickets {
            shortDescriptions[i] = ticket.ShortDescription 
        }
        var err error
        clusters, err = database.TFIDFKMeansClustering(shortDescriptions)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        h.Cache.mu.Lock()
        h.Cache.lastIncidents = *incidents
        h.Cache.lastClusters = clusters
        h.Cache.lastUpdateTime = time.Now()
        h.Cache.mu.Unlock()
    }

    response := createClusteredTicketResponse(clusters, ToTickets(incidents))
    jsonResponse(w, response)
}

func createClusteredTicketResponse(clusters database.TicketResponse, tickets []models.Ticket) ClusteredTicketResponse {
    response := ClusteredTicketResponse{
        Clusters: make([]ClusteredTickets, len(clusters.Clusters)),
    }

    ticketMap := make(map[string]models.Ticket)
    for _, ticket := range tickets {
        ticketMap[ticket.ShortDescription] = ticket
    }

    for i, cluster := range clusters.Clusters {
        clusteredTickets := ClusteredTickets{
            ClusterDescription: cluster.ClusterDescription,
            Tickets:            make([]models.Ticket, 0),
        }

        for _, textEntry := range cluster.TextEntries {
            if ticket, found := ticketMap[textEntry]; found {
                clusteredTickets.Tickets = append(clusteredTickets.Tickets, ticket)
            }
        }

        response.Clusters[i] = clusteredTickets
    }

    return response
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
