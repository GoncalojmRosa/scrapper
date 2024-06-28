package store

import (
	"log"
	"testing"
)

func TestRedisDB(t *testing.T) {
	rdb, err := NewRedisStorage()
	if err != nil {
		t.Fatalf("NewRedisStorage() err = %v; want nil", err)
	}
	log.Printf("Redis connection established: %v", rdb)
}

func TestInsertProduct(t *testing.T) {
	rdb, err := NewRedisStorage()
	if err != nil {
		log.Fatalf("NewRedisStorage() err = %v; want nil", err)
	}
	store := NewStore(rdb)
	err = store.InsertProduct("test", "test", "test")
	if err != nil {
		log.Fatalf("InsertProduct() err = %v; want nil", err)
	}
}
