package store

import (
	"github.com/GoncalojmRosa/scrapper/types"
	"github.com/go-redis/redis"
)

type Storage struct {
	rdb *redis.Client
}

type Store interface {
	InsertProduct(name, price, img string) error
	GetProducts() ([]types.Product, error)
}

func NewStore(rdb *redis.Client) *Storage {
	return &Storage{
		rdb: rdb,
	}
}

func (s *Storage) InsertProduct(name, price, img string) error {
	return s.rdb.HSet("products", name, price).Err()
}
