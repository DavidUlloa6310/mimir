package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type CuMLClusterResult struct {
	Clusters []struct {
		ClusterDescription string   `json:"cluster_description"`
		TextEntries        []string `json:"text_entries"`
	} `json:"clusters"`
}

func CuMLTFIDFKMeansClustering(documents []string) (TicketResponse, error) {
	documentsJSON, err := json.Marshal(documents)
	if err != nil {
		log.Printf("Error marshaling documents to JSON: %v", err)
		return TicketResponse{}, fmt.Errorf("error marshaling documents: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", "-c", `
import json
import sys
sys.path.append('./backend')  # Adjust path as needed
from database.clustering_cuml import tfidf_kmeans_clustering_api

# Read input from stdin
documents = json.loads(sys.stdin.read())
result = tfidf_kmeans_clustering_api(documents)
print(result)
`)

	cmd.Stdin = strings.NewReader(string(documentsJSON))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing Python script: %v, Output: %s", err, output)

		log.Printf("Falling back to Go implementation")
		return TFIDFKMeansClustering(documents)
	}

	var result CuMLClusterResult
	if err := json.Unmarshal(output, &result); err != nil {
		log.Printf("Error parsing Python output: %v, Output: %s", err, output)
		return TicketResponse{}, fmt.Errorf("error parsing Python output: %v", err)
	}

	response := TicketResponse{
		Clusters: make([]ClusterEntry, len(result.Clusters)),
	}

	for i, cluster := range result.Clusters {
		response.Clusters[i] = ClusterEntry{
			ClusterDescription: cluster.ClusterDescription,
			TextEntries:        cluster.TextEntries,
		}
	}

	return response, nil
}

func TFIDFKMeansClusteringWithFallback(documents []string) (TicketResponse, error) {
	response, err := CuMLTFIDFKMeansClustering(documents)
	if err != nil {
		log.Printf("cuML clustering failed, falling back to Go implementation: %v", err)
		return TFIDFKMeansClustering(documents)
	}
	return response, nil
}
