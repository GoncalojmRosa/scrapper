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
		product.Price = e.ChildText(".ct-price-formatted")

		// adding the product instance with scraped data to the list of products
		products = append(products, product)
	})

	col.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	})

	col.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Visit the website after setting up all the callbacks
	col.Visit("https://www.continente.pt/campanhas/campanhas-folhetos/folheto-semanal-2/?start=0&srule=Trading-categorias-destaques&pmin=0.01")

	// Print the products after the visit is complete
	fmt.Println(products)
	fmt.Println("PRODUCT LENGTH:", len(products))
}
