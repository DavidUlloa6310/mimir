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
    suggestionsHandler := handlers.NewSuggestionsHandler()
    chatHandler := handlers.NewChatHandler()
    docHandler := handlers.NewDocumentationHandler()
	authHandler := handlers.NewAuthorizationHandler()

    http.Handle("/tickets", handlers.AuthMiddleware(http.HandlerFunc(ticketHandler.TicketsHandler)))
    http.Handle("/suggestions", handlers.AuthMiddleware(http.HandlerFunc(suggestionsHandler.SuggestionsHandler)))
    http.Handle("/chat", handlers.AuthMiddleware(http.HandlerFunc(chatHandler.ChatHandler)))
    http.Handle("/documentation", handlers.AuthMiddleware(http.HandlerFunc(docHandler.DocumentationHandler)))
	http.Handle("/authorization", http.HandlerFunc(authHandler.AuthorizationHandler))

    fmt.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
