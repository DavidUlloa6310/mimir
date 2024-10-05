package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/davidulloa/mimir/database"
)

type TicketRequestBody struct {
    InstanceID string `json:"instance_id"`
}

func ParseCredentials(r *http.Request) (string, string, string, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return "",  "", "", errors.New("basic authentication could not be collected from request")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", "", "", err
	}

	var responseBody TicketRequestBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", "", "", err
	}

	if responseBody.InstanceID == "" {
		return "", "", "", errors.New("`instance_id` not passed into request body")
	}

	return responseBody.InstanceID, username, password, nil
}

type AuthorizationHandler struct {}

// NewDocumentationHandler creates and returns a new DocumentationHandler instance
func NewAuthorizationHandler() *AuthorizationHandler{
	return &AuthorizationHandler{}
}

func AuthMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instanceID, username, password, err := ParseCredentials(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		valid, err := database.ValidateAuthentication(instanceID, username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !valid {
			http.Error(w, "authentication with the given username nad password could not validated", http.StatusBadRequest)
			return
		}

        handler.ServeHTTP(w, r)
    })
}

// DocumentationHandler serves the documentation route and returns dummy data
func (h *AuthorizationHandler) AuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	instanceID, username, password, err := ParseCredentials(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.RegisterAuthentication(instanceID, username, password)
	if err != nil {
		errMsg := fmt.Sprintf("Could not validate authentication: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authorization granted"))
}