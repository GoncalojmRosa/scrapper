package main

import "github.com/GoncalojmRosa/scrapper/handlers"

func main() {
	router := NewRouter()
	handler := handlers.()

	router.HandleFunc("/api/v1/health", handler.Health).Methods("GET")
}
