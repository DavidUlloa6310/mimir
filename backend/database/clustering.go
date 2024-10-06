package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"

	"math"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

type TFIDFVectorizer struct {
    Vocabulary map[string]int
    IDF        map[string]float64
}

func NewTFIDFVectorizer() *TFIDFVectorizer {
    return &TFIDFVectorizer{
        Vocabulary: make(map[string]int),
        IDF:        make(map[string]float64),
    }
}

type ClusterEntry struct {
	ClusterDescription string   `json:"cluster_description"`
	TextEntries        []string `json:"text_entries"`
}

type TicketResponse struct {
	Clusters []ClusterEntry `json:"clusters"`
}


func (v *TFIDFVectorizer) FitTransform(descriptions []string) [][]float64 {
    docFreq := make(map[string]int)
    for _, desc := range descriptions {
        tokens := strings.Fields(strings.ToLower(desc))
        uniqueTokens := make(map[string]bool)
        for _, token := range tokens {
            if _, exists := v.Vocabulary[token]; !exists {
                v.Vocabulary[token] = len(v.Vocabulary)
            }
            uniqueTokens[token] = true
        }
        for token := range uniqueTokens {
            docFreq[token]++
        }
    }

    numDocs := float64(len(descriptions))
    for token, freq := range docFreq {
        v.IDF[token] = math.Log(numDocs / float64(freq))
    }

    vectors := make([][]float64, len(descriptions))
    for i, desc := range descriptions {
        vector := make([]float64, len(v.Vocabulary))
        tokens := strings.Fields(strings.ToLower(desc))
        for _, token := range tokens {
            if idx, exists := v.Vocabulary[token]; exists {
                tf := float64(strings.Count(desc, token)) / float64(len(tokens))
                vector[idx] = tf * v.IDF[token]
            }
        }
        vectors[i] = vector
    }

    return vectors
}

func coordinatesEqual(a, b clusters.Coordinates) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func clusterTexts(tfidfMatrix [][]float64, numClusters int) map[int][]int {
    var observations clusters.Observations
    for _, vector := range tfidfMatrix {
        observations = append(observations, clusters.Coordinates(vector))
    }

    km, err := kmeans.NewWithOptions(0.01, nil)
    if err != nil {
        panic(err)
    }

    clusters, err := km.Partition(observations, numClusters)
    if err != nil {
        panic(err)
    }

    clusteredTexts := make(map[int][]int)
    for clusterIndex, cluster := range clusters {
        for _, observation := range cluster.Observations {
            for i, orig := range observations {
                if coordinatesEqual(orig.Coordinates(), observation.Coordinates()) {
                    clusteredTexts[clusterIndex] = append(clusteredTexts[clusterIndex], i)
                    break
                }
            }
        }
    }

    return clusteredTexts
}


func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var TicketResponseSchema = GenerateSchema[TicketResponse]()
func generateTicketDescriptions(clusters [][]string) (*TicketResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	promptEngineering := "You are a helpful assistant. Can you provide a short summary of the main features of these products? Write a very short (5 words maximum title that summarizes all of them collectively). Usually, you would include the product name that most correlates to those incident reports."

	clustersJSON, err := json.Marshal(clusters)
	if err != nil {
		return nil, fmt.Errorf("error marshaling clusters: %v", err)
	}

	content := fmt.Sprintf("Classify the following clusters:\n%s", string(clustersJSON))

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("ticket_response"),
		Description: openai.F("Clustered ticket descriptions"),
		Schema:      openai.F(TicketResponseSchema),
		Strict:      openai.Bool(true),
	}

	chat, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(promptEngineering),
			openai.UserMessage(content),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})

	if err != nil {
		return nil, fmt.Errorf("error getting chat completion: %v", err)
	}

	if len(chat.Choices) == 0 {
		return nil, fmt.Errorf("no response from the model")
	}

	var response TicketResponse
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return &response, nil
}


func TFIDFKMeansClustering(documents []string) (TicketResponse, error) {
    vectorizer := NewTFIDFVectorizer()
    tfidfMatrix := vectorizer.FitTransform(documents)

    numClusters := 3
    clusteredTexts := clusterTexts(tfidfMatrix, numClusters)

    clusters := make([][]string, numClusters)
    for clusterIndex, textIndices := range clusteredTexts {
        for _, textIndex := range textIndices {
            clusters[clusterIndex] = append(clusters[clusterIndex], documents[textIndex])
        }
    }

    response, err := generateTicketDescriptions(clusters)
    if err != nil {
        return TicketResponse{}, fmt.Errorf("Error generating ticket descriptions: %v", err) 
    }

    if len(response.Clusters) != numClusters {
        fmt.Printf("Expected %d clusters, got %d", numClusters, len(response.Clusters))
    }

    return *response, nil 
}
