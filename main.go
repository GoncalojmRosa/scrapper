package main

import (
	"github.com/GoncalojmRosa/scrapper/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	handler := handlers.New()

	router.HandleFunc("/api/v1/health", handler.HandleHome).Methods("GET")
}
