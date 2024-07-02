package store

import (
	"github.com/GoncalojmRosa/scrapper/types"
	"github.com/go-redis/redis"
)

type Storage struct {
	rdb *redis.Client
}

type Store interface {
	GetProducts() ([]types.Product, error)
}

func NewStore(rdb *redis.Client) *Storage {
	return &Storage{
		rdb: rdb,
	}
}

func (s *Storage) GetProducts() ([]types.Product, error) {
	products := []types.Product{}
	result, err := s.rdb.HGetAll("auchan").Result()
	if err != nil {
		return nil, err
	}
	for name, price := range result {
		products = append(products, types.Product{
			Name:  name,
			Price: price,
		})
	}
	return products, nil
}
