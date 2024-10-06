package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/davidulloa/mimir/database"
	"github.com/davidulloa/mimir/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// SuggestionsHandler handles suggestions-related requests
type SuggestionsHandler struct {
	// You can add dependencies here, such as database connections, services, etc.
}

type SuggestionsBody struct {
	ticketIds []string `json:"tickets"`
}

func NewSuggestionsHandler() *SuggestionsHandler {
	return &SuggestionsHandler{}
}

type SuggestionOpenAiSchema struct {
	suggestion []Suggestion
}

type Suggestion struct {
	description string
	accelerator models.Accelerator
}

func GenerateSuggestions(clusters []database.ClusterEntry, accelerators []models.Accelerator) (SuggestionOpenAiSchema, error) {

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)

	suggestionPrompt := fmt.Sprintf(`# Task
	You're a ServiceNow Accelerator Assistant - meant to find the best support for ServiceNow users in ServiceNow accelerators.
	By understanding what their needs are as customers and the accelerators which are available, you can find the best possible path for customers to better user SeriveNow products.

	## Inputs
	You are given the accelerators from the ServiceNow catalog and the clusters of a user's most relevant issues. Use these pieces of information to describe potential solutions and list the associated accelerator.\n
	Clusters:
	%s
	Accelerators:
	%s
	`, clusters, accelerators)


	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("suggestion_response"),
		Description: openai.F("Provides practical suggestion tied to accelerator based on ticket used"),
		Schema:      openai.F(database.GenerateSchema[SuggestionOpenAiSchema]()),
		Strict:      openai.Bool(true),
	}

	chat, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(suggestionPrompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	 })

	var response SuggestionOpenAiSchema
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &response)
	if err != nil {
		return SuggestionOpenAiSchema{}, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return response, nil
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
	var descriptions []string
	for _, ticket := range tickets {
		descriptions = append(descriptions, ticket.ShortDescription)
	}

	clusters, err := database.TFIDFKMeansClustering(descriptions)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accelerators, err := database.GetAllAccelerators()
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	suggestions, err := GenerateSuggestions(clusters.Clusters, accelerators)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}