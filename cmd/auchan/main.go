package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gocolly/colly"
)

type AuchanResponse struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Img   string `json:"img"`
}

var redisClient *redis.Client

func ConnectToRedis() {
	opt, err := redis.ParseURL("")
	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(opt)
	fmt.Println(redisClient.Ping())
}

func SaveToRedis(key string, value string) {
	// TODO: Save other values, like Img, instead only price

	err := redisClient.HSet("auchan", key, value).Err()
	if err != nil {
		fmt.Println("Error saving to redis:", err)
		return
	}
}

// var proxies = []string{
// 	"https://95.92.206.174:3128",
// 	"https://188.93.237.29:3128",
// }

func main() {
	ConnectToRedis()

	col := colly.NewCollector()
	col.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// col.SetRequestTimeout(120 * time.Second)

	// col.SetProxyFunc(func(_ *http.Request) (*url.URL, error) {
	// 	proxyURL, _ := url.Parse(proxies[rand.Intn(len(proxies))])
	// 	return proxyURL, nil
	// })
	// iterating over the list of HTML product elements

	products := make([]AuchanResponse, 0)

	col.OnHTML("div.product", func(e *colly.HTMLElement) {
		// initializing a new Product instance
		product := AuchanResponse{}

		result := e.ChildAttr("div.product-tile", "data-gtm")
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
		product.Img = e.ChildAttr("link", "href")

		// adding the product instance with scraped data to the list of products
		products = append(products, product)
	})

	col.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	col.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Visit the website after setting up all the callbacks
	col.Visit("https://www.auchan.pt/on/demandware.store/Sites-AuchanPT-Site/pt_PT/Search-UpdateGrid?cgid=produtos-frescos&prefn1=soldInStores&prefv1=000&start=0&sz=200")
	col.Wait()

	fmt.Println("PRODUCT LENGTH:", len(products))
	for _, product := range products {
		SaveToRedis(product.Name, product.Price)
	}
}
