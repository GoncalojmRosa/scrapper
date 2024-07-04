package main

import (
	"log"
	"net/http"

	"github.com/GoncalojmRosa/scrapper/handlers"
	"github.com/GoncalojmRosa/scrapper/store"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main() {
	rdb, err := store.NewRedisStorage(store.Envs.RedisURL)

	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStore(rdb)

	initStorage(rdb)

	router := mux.NewRouter()

	handler := handlers.New(store)

	router.HandleFunc("/", handler.HandleListProducts).Methods("GET")
	router.HandleFunc("/products", handler.HandleProductSearch).Methods("POST")

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.ListenAndServe(":8080", router)
}

func initStorage(rdb *redis.Client) {
	err := rdb.Ping().Err()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("Redis connection established")
}
