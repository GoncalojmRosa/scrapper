package store

import (
	"fmt"
	"log"

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
	// Get all product names from the set
	productNames, err := s.rdb.SMembers("auchan:products").Result()
	if err != nil {
		log.Fatal("Error retrieving product names:", err)
		return nil, err
	}

	var products []types.Product
	for _, name := range productNames {
		key := fmt.Sprintf("auchan:products:%s", name)
		data, err := s.rdb.HGetAll(key).Result()
		if err != nil {
			log.Println("Error retrieving product data:", err)
			continue
		}

		// Ensure the required fields are present in the data
		name, nameOk := data["name"]
		price, priceOk := data["price"]
		img, imgOk := data["img"]

		if !nameOk || !priceOk || !imgOk {
			log.Printf("Missing data for product: %s\n", name)
			continue
		}

		productData := types.Product{
			Name:  name,
			Price: price,
			Img:   img,
		}

		products = append(products, productData)
	}

	return products, nil
}
