package main

import (
	"net/http"

	"github.com/GoncalojmRosa/scrapper/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	handler := handlers.New()

	router.HandleFunc("/", handler.HandleHome).Methods("GET")
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.ListenAndServe(":8080", router)
}
