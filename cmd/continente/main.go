package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/gocolly/colly"
)

type ContinenteResponse struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Img   string  `json:"img"`
}

var products []ContinenteResponse

var redisClient *redis.Client

func ConnectToRedis() {
	opt, err := redis.ParseURL("")
	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(opt)
	fmt.Println(redisClient.Ping())
}

func SaveToRedis(product ContinenteResponse) {
	key := fmt.Sprintf("products:%s", strings.Split(strings.ToUpper(product.Name), " ")[0])

	// Save individual fields to Redis hash
	err := redisClient.HSet(key, "name", product.Name).Err()
	if err != nil {
		fmt.Println("Error saving name to redis:", err)
		return
	}

	err = redisClient.HSet(key, "price", product.Price).Err()
	if err != nil {
		fmt.Println("Error saving price to redis:", err)
		return
	}

	err = redisClient.HSet(key, "img", product.Img).Err()
	if err != nil {
		fmt.Println("Error saving img to redis:", err)
		return
	}

	err = redisClient.HSet(key, "supermarket", "continente").Err()
	if err != nil {
		fmt.Println("Error saving supermarket to redis:", err)
		return
	}

	// Add product name to the set of product names
	err = redisClient.SAdd("products", product.Name).Err()
	if err != nil {
		fmt.Println("Error adding product to set:", err)
		return
	}

	// Set key to expire in one week
	err = redisClient.Expire(key, 7*24*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting key expiration:", err)
	}
}

func main() {
	ConnectToRedis()
	col := colly.NewCollector(colly.Async(true))
	col.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	col.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	// iterating over the list of HTML product elements
	col.OnHTML("div.product", func(e *colly.HTMLElement) {
		product := ContinenteResponse{}
		result := e.ChildAttr("div.product-tile", "data-product-tile-impression")
		if !json.Valid([]byte(result)) {
			fmt.Println("Invalid JSON, skipping:", result)
			return
		}
		err := json.Unmarshal([]byte(result), &product)
		if err != nil {
			fmt.Println(result)
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		product.Img = e.ChildAttr("img.ct-tile-image", "data-src")
		products = append(products, product)
		SaveToRedis(product)
	})

	col.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	})

	col.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Function to visit each page and collect products
	visitPage := func(pageNumber int) error {
		url := fmt.Sprintf("https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?&start=%d&sz=24", pageNumber*24)
		return col.Visit(url)
	}

	// Iterate over pages until no more products are found
	pageNumber := 0
	for {
		err := visitPage(pageNumber)
		if err != nil {
			fmt.Println("Error visiting page:", err)
			break
		}

		// Wait for requests to finish before deciding to stop or continue
		col.Wait()

		if len(products) <= pageNumber*24 {
			break
		}
		pageNumber++
	}

	fmt.Println("PRODUCT LENGTH:", len(products))
}
