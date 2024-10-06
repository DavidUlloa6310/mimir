package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/davidulloa/mimir/database"
	"github.com/davidulloa/mimir/models"
)

type DocumentationHandler struct {}

func NewDocumentationHandler() *DocumentationHandler {
	return &DocumentationHandler{}
}

func (h *DocumentationHandler) verifyBasicAuth(r *http.Request, instanceID string) error {
    username, password, ok := r.BasicAuth()
    if !ok {
        return fmt.Errorf("basic authentication required")
    }

    validated, err := database.ValidateAuthentication(instanceID, username, password)
    if err != nil {
        return fmt.Errorf("authentication validation error: %v", err)
    }

    if !validated {
        return fmt.Errorf("invalid credentials")
    }

    return nil
}

func (h *DocumentationHandler) DocumentationHandler(w http.ResponseWriter, r *http.Request) {
    var requestBody struct {
        InstanceID     string   `json:"instanceId"`
        AcceleratorIds []string `json:"acceleratorIds"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if requestBody.InstanceID == "" {
        http.Error(w, "instanceId is required", http.StatusBadRequest)
        return
    }

    if err := h.verifyBasicAuth(r, requestBody.InstanceID); err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    documentation, err := h.scrapeDocumentations(requestBody.AcceleratorIds)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(documentation); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (h *DocumentationHandler) scrapeDocumentations(ids []string) ([]models.Documentation, error) {
	var wg sync.WaitGroup
	documentation := make([]models.Documentation, 0, len(ids))
	errs := make(chan error, len(ids))

	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			accelerator, err := database.GetAcceleratorByID(id)
			if err != nil {
				errs <- err
				return
			}

			title, err := h.scrapeTitle(accelerator.Url)
			if err != nil {
				errs <- err
				return
			}

			doc := models.Documentation{
				Title:      title,
				Accelerator: *accelerator,
			}
			documentation = append(documentation, doc)
		}(id)
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return nil, <-errs 
	}

	return documentation, nil
}

func (h *DocumentationHandler) scrapeTitle(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse document: %w", err)
	}

	title := doc.Find("title").Text()
	return title, nil
}