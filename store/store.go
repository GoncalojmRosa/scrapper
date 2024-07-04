package store

import (
	"fmt"
	"log"
	"strings"

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
	productNames, err := s.rdb.SMembers("products").Result()
	if err != nil {
		log.Fatal("Error retrieving product names:", err)
		return nil, err
	}

	var products []types.Product
	counter := 0
	limit := 20

	for _, name := range productNames {
		if counter >= limit {
			break
		}

		key := fmt.Sprintf("products:%s", strings.Split(name, " ")[0])
		data, err := s.rdb.HGetAll(key).Result()
		if err != nil {
			log.Println("Error retrieving product data:", err)
			continue
		}

		// Ensure the required fields are present in the data
		name, nameOk := data["name"]
		price, priceOk := data["price"]
		img, imgOk := data["img"]
		sp, spOk := data["supermarket"]

		if !nameOk || !priceOk || !imgOk || !spOk {
			log.Printf("Missing data for product: %s\n", name)
			continue
		}

		productData := types.Product{
			Name:        name,
			Price:       price,
			Img:         img,
			Supermarket: sp,
		}

		products = append(products, productData)
		counter++
	}

	return products, nil
}

func (s *Storage) GetProductByName(name string) ([]types.Product, error) {
	product, err := s.rdb.HGetAll("products:" + strings.ToUpper(name)).Result()
	if err != nil {
		log.Println("Error retrieving product data:", err)
		return nil, err
	}
	// Ensure the required fields are present in the data
	name, nameOk := product["name"]
	price, priceOk := product["price"]
	img, imgOk := product["img"]
	sp, spOk := product["supermarket"]

	if !nameOk || !priceOk || !imgOk || !spOk {
		log.Printf("Missing data for product: %s\n", name)
		return nil, fmt.Errorf("missing data for product: %s", name)
	}

	productData := types.Product{
		Name:        name,
		Price:       price,
		Img:         img,
		Supermarket: sp,
	}

	return []types.Product{productData}, nil
}
