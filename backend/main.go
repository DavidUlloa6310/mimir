package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidulloa/mimir/handlers"
)

func main() {
	client := &http.Client{}

	ticketHandler := handlers.NewTicketHandler(client)
	http.HandleFunc("/tickets", ticketHandler.TicketsHandler)

	suggestionsHandler := handlers.NewSuggestionsHandler()
	http.HandleFunc("/suggestions", suggestionsHandler.SuggestionsHandler)

	chatHandler := handlers.NewChatHandler()
	http.HandleFunc("/chat", chatHandler.ChatHandler)

	docHandler := handlers.NewDocumentationHandler()
	http.HandleFunc("/documentation", docHandler.DocumentationHandler)

	// Start the server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
