package main

import (
	"encoding/json"
	"fmt"

	"github.com/gocolly/colly"
)

type ContinenteResponse struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Img   string  `json:"img"`
}

var products []ContinenteResponse

func main() {
	col := colly.NewCollector()
	col.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// iterating over the list of HTML product elements
	col.OnHTML("div.product", func(e *colly.HTMLElement) {
		// initializing a new Product instance
		product := ContinenteResponse{}

		// scraping the data of interest
		result := e.ChildAttr("div.product-tile", "data-product-tile-impression")
		err := json.Unmarshal([]byte(result), &product)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
		}
		product.Img = e.ChildAttr("img.ct-tile-image", "data-src")
		products = append(products, product)
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
		url := fmt.Sprintf("https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=mercearias&pmin=0%%2e01&srule=FOOD_Mercearia&start=%d&sz=24", pageNumber*24)
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

	fmt.Println("Products:", products)
	fmt.Println("PRODUCT LENGTH:", len(products))
}
