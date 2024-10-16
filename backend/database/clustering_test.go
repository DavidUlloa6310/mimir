package database

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/muesli/clusters"
)

func TestGenerateTicketDescriptions(t *testing.T) {
  
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Fatal("OPENAI_API_KEY is not set in the environment")
	}

	clusters := [][]string{
		{"apple", "banana", "orange"},
		{"dog", "cat", "hamster"},
		{"car", "bus", "train"},
	}

	response, err := generateTicketDescriptions(clusters)
	if err != nil {
		t.Fatalf("Error generating ticket descriptions: %v", err)
	}

	fmt.Printf("Structured Response:\n")
	for i, cluster := range response.Clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("  Description: %s\n", cluster.ClusterDescription)
		fmt.Printf("  Text Entries:\n")
		for _, entry := range cluster.TextEntries {
			fmt.Printf("    - %s\n", entry)
		}
		fmt.Println()
	}

	if len(response.Clusters) != len(clusters) {
		t.Errorf("Expected %d clusters, but got %d", len(clusters), len(response.Clusters))
	}

	for i, cluster := range response.Clusters {
		if len(cluster.TextEntries) != len(clusters[i]) {
			t.Errorf("Cluster %d: Expected %d entries, but got %d", i+1, len(clusters[i]), len(cluster.TextEntries))
		}
		if cluster.ClusterDescription == "" {
			t.Errorf("Cluster %d: Description is empty", i+1)
		}
	}
}

func TestNewTFIDFVectorizer(t *testing.T) {
	vectorizer := NewTFIDFVectorizer()

	if vectorizer == nil {
		t.Fatal("NewTFIDFVectorizer returned nil")
	}

	if vectorizer.Vocabulary == nil {
		t.Error("Vocabulary map is nil")
	}

	if vectorizer.IDF == nil {
		t.Error("IDF map is nil")
	}

	if len(vectorizer.Vocabulary) != 0 {
		t.Errorf("Expected empty Vocabulary, got %d items", len(vectorizer.Vocabulary))
	}

	if len(vectorizer.IDF) != 0 {
		t.Errorf("Expected empty IDF, got %d items", len(vectorizer.IDF))
	}
}

func TestTFIDFVectorizerFitTransform(t *testing.T) {
	vectorizer := NewTFIDFVectorizer()
	descriptions := []string{
		"This is a test",
		"This is another test",
		"And this is a third test",
	}

	vectors := vectorizer.FitTransform(descriptions)

	if len(vectors) != len(descriptions) {
		t.Errorf("Expected %d vectors, got %d", len(descriptions), len(vectors))
	}

	expectedVocabSize := 7 // "this", "is", "a", "test", "another", "and", "third"
	if len(vectorizer.Vocabulary) != expectedVocabSize {
		t.Errorf("Expected vocabulary size of %d, got %d", expectedVocabSize, len(vectorizer.Vocabulary))
	}

	if len(vectorizer.IDF) != expectedVocabSize {
		t.Errorf("Expected IDF size of %d, got %d", expectedVocabSize, len(vectorizer.IDF))
	}

	for i, vector := range vectors {
		if len(vector) != expectedVocabSize {
			t.Errorf("Vector %d: expected length %d, got %d", i, expectedVocabSize, len(vector))
		}
	}

	fmt.Println("Vocabulary:", vectorizer.Vocabulary)
}

func TestCoordinatesEqual(t *testing.T) {
	tests := []struct {
		a, b     clusters.Coordinates
		expected bool
	}{
		{clusters.Coordinates{1, 2, 3}, clusters.Coordinates{1, 2, 3}, true},
		{clusters.Coordinates{1, 2, 3}, clusters.Coordinates{1, 2, 4}, false},
		{clusters.Coordinates{1, 2, 3}, clusters.Coordinates{1, 2}, false},
		{clusters.Coordinates{}, clusters.Coordinates{}, true},
	}

	for _, test := range tests {
		result := coordinatesEqual(test.a, test.b)
		if result != test.expected {
			t.Errorf("coordinatesEqual(%v, %v) = %v; want %v", test.a, test.b, result, test.expected)
		}
	}
}

func TestClusterTexts(t *testing.T) {
	tfidfMatrix := [][]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 0},
		{0, 1, 1},
	}
	numClusters := 3

	result := clusterTexts(tfidfMatrix, numClusters)

	if len(result) != numClusters {
		t.Errorf("Expected %d clusters, got %d", numClusters, len(result))
	}

	totalTexts := 0
	for _, texts := range result {
		totalTexts += len(texts)
	}

	if totalTexts != len(tfidfMatrix) {
		t.Errorf("Expected %d total texts, got %d", len(tfidfMatrix), totalTexts)
	}

	// Check if all text indices are unique and within range
	allIndices := make(map[int]bool)
	for _, texts := range result {
		for _, idx := range texts {
			if idx < 0 || idx >= len(tfidfMatrix) {
				t.Errorf("Invalid text index: %d", idx)
			}
			if allIndices[idx] {
				t.Errorf("Duplicate text index: %d", idx)
			}
			allIndices[idx] = true
		}
	}
}

func TestClusterEntry(t *testing.T) {
	entry := ClusterEntry{
		ClusterDescription: "Test Cluster",
		TextEntries:        []string{"Entry 1", "Entry 2", "Entry 3"},
	}

	if entry.ClusterDescription != "Test Cluster" {
		t.Errorf("Expected ClusterDescription 'Test Cluster', got '%s'", entry.ClusterDescription)
	}

	if len(entry.TextEntries) != 3 {
		t.Errorf("Expected 3 TextEntries, got %d", len(entry.TextEntries))
	}

	expectedEntries := []string{"Entry 1", "Entry 2", "Entry 3"}
	if !reflect.DeepEqual(entry.TextEntries, expectedEntries) {
		t.Errorf("TextEntries mismatch. Expected %v, got %v", expectedEntries, entry.TextEntries)
	}
}

func TestTicketResponse(t *testing.T) {
	response := TicketResponse{
		Clusters: []ClusterEntry{
			{
				ClusterDescription: "Cluster 1",
				TextEntries:        []string{"Entry 1", "Entry 2"},
			},
			{
				ClusterDescription: "Cluster 2",
				TextEntries:        []string{"Entry 3", "Entry 4", "Entry 5"},
			},
		},
	}

	if len(response.Clusters) != 2 {
		t.Errorf("Expected 2 Clusters, got %d", len(response.Clusters))
	}

	if response.Clusters[0].ClusterDescription != "Cluster 1" {
		t.Errorf("Expected ClusterDescription 'Cluster 1', got '%s'", response.Clusters[0].ClusterDescription)
	}

	if len(response.Clusters[1].TextEntries) != 3 {
		t.Errorf("Expected 3 TextEntries in Cluster 2, got %d", len(response.Clusters[1].TextEntries))
	}
}

func TestTFIDFKMeansClustering(t *testing.T) {
	documents := []string{
		"Jira wasn't working for me this morning, couldn't log in after several tries.",
		"Confluence kept crashing every time I tried to save a page, super frustrating.",
		"Jira tickets are not syncing with our Slack channel since yesterday.",
		"Couldn't upload any attachments to Jira, files were stuck at 0%.",
		"Jira board was super slow, it took minutes just to open one ticket.",
		"Confluence search feature has been down, unable to find important documents.",
		"Trello wasn't displaying my boards, got a blank screen instead.",
		"Jira mobile app froze as I tried to update a ticket, had to restart the app.",
		"Jira notifications are delayed by hours, I'm missing important updates.",
		"Jira sprint planning feature is glitchy, we lost all story points for a session.",
		"Trello cards are not moving to different columns when dragged.",
		"Jira API is throwing 500 errors on every request for ticket data.",
		"Can't close or transition Jira tickets, error message pops up every time.",
		"Confluence editor keeps refreshing randomly, making it impossible to edit pages.",
		"Jira filters aren't working properly, not returning correct ticket results.",
		"Confluence permissions seem off, some users can't access shared spaces.",
		"Jira webhooks aren't triggering, automation rules are not working as expected.",
		"Trello power-ups aren't loading, integrations seem broken.",
		"Jira comments are disappearing after refreshing the page.",
		"Can't export reports from Jira, the export button is completely unresponsive.",
		"Jira backlog view is completely empty, even though tickets are assigned.",
		"Confluence page versions not showing up, can't revert to previous versions.",
	}

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
		t.Fatalf("Error generating ticket descriptions: %v", err)
	}

	if len(response.Clusters) != numClusters {
		t.Errorf("Expected %d clusters, got %d", numClusters, len(response.Clusters))
	}

	for i, cluster := range response.Clusters {
		if len(cluster.ClusterDescription) == 0 {
			t.Errorf("Cluster %d description is empty", i)
		}
		if len(cluster.TextEntries) != len(clusteredTexts[i]) {
			t.Errorf("Cluster %d: expected %d text entries, got %d", i, len(clusteredTexts[i]), len(cluster.TextEntries))
		}
	}

	fmt.Println("Clustering Results:")
	for i, cluster := range response.Clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("  Description: %s\n", cluster.ClusterDescription)
		fmt.Printf("  Text Entries:\n")
		for _, entry := range cluster.TextEntries {
			fmt.Printf("    - %s\n", entry)
		}
		fmt.Println()
	}
}
