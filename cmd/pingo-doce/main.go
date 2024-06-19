package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Product struct {
	Source Source `json:"_source"`
}

type Source struct {
	FirstName string  `json:"firstName"`
	Price     float64 `json:"regularPrice"` // Assuming the price field exists
}

type Section struct {
	Products []Product `json:"products"`
}

type Response struct {
	Sections map[string]Section `json:"sections"`
}

func main() {
	url := "https://mercadao.pt/api/Campaigns/products/catalogue/6107d28d72939a003ff6bf51?from=0&brandIds=%5B%5D&categoryIds=%5B%5D&mainCategoryIds=%5B%22621384a88c69c5003ff5cf5e%22%5D&sort=%7B%22activePromotion%22:%22desc%22%7D&ids=%5B%22621384a88c69c5003ff5cf5e%22%5D&size=100"

	// Create an HTTP client
	client := &http.Client{}

	// Create HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set user-agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	// Make HTTP GET request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making GET request:", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected HTTP status code: %d", resp.StatusCode)
	}

	// Decode JSON response
	var product Response
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	log.Println(product)
}
