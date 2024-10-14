package main

import (
    "os"
	"fmt"
	"log"
	"net/http"

	"github.com/davidulloa/mimir/handlers"
)

func enableCORS(next http.Handler) http.Handler {
    frontend := os.Getenv("FRONTEND_IP")
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set the necessary headers
        w.Header().Set("Access-Control-Allow-Origin", frontend)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // If it's an OPTIONS request, end here
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Proceed with the next handler
        next.ServeHTTP(w, r)
    })
}


func main() {
	client := &http.Client{}

    ticketHandler := handlers.NewTicketHandler(client)
    suggestionsHandler := handlers.NewSuggestionsHandler()
    chatHandler := handlers.NewChatHandler()
    docHandler := handlers.NewDocumentationHandler()
	authHandler := handlers.NewAuthorizationHandler()

    http.Handle("/tickets", enableCORS(handlers.AuthMiddleware(http.HandlerFunc(ticketHandler.TicketsHandler))))
    http.Handle("/suggestions", enableCORS(handlers.AuthMiddleware(http.HandlerFunc(suggestionsHandler.SuggestionsHandler))))
    http.Handle("/chat", enableCORS(handlers.AuthMiddleware(http.HandlerFunc(chatHandler.ChatHandler))))
    http.Handle("/documentation", enableCORS(handlers.AuthMiddleware(http.HandlerFunc(docHandler.DocumentationHandler))))
	http.Handle("/authorization", enableCORS(http.HandlerFunc(authHandler.AuthorizationHandler)))

    fmt.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
