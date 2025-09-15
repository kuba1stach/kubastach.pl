package main

import (
	"log"
	"net/http"

	"kubastach.pl/backend/pkg/api"
)

func main() {
	// Create a new instance of our server implementation.
	server := api.NewServer()

	// Create a new router.
	router := api.NewRouter(server)

	// Start the server.
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
