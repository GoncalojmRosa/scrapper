package main

import (
	"fmt"

	"github.com/GoncalojmRosa/scrapper/types"
	"github.com/gocolly/colly"
)

var products []types.Product

func main() {
	col := colly.NewCollector()
	col.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// iterating over the list of HTML product elements
	col.OnHTML("div.product", func(e *colly.HTMLElement) {
		// initializing a new Product instance
		product := types.Product{}

		// scraping the data of interest
		product.Img = e.ChildAttr("img", "src")
		product.Name = e.ChildText("h2")
		product.Price = e.ChildText(".value")

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
	col.Visit("https://www.continente.pt/")

	// Print the products after the visit is complete
	fmt.Println(products)
}
