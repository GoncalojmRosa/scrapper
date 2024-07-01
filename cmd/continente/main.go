package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type ContinenteResponse struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Img   string  `json:"img"`
}

var (
	urlsToScrape = []string{
		"https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=biologicos&pmin=0%2e01&start=",
		"https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=laticinios&pmin=0%2e01&start=",
		"https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=peixaria-e-talho&pmin=0%2e01&start=",
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=frutas-legumes&pmin=0%2e01&start=",
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=charcutaria-queijo&pmin=0%2e01&start=",
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=padaria-e-pastelaria&pmin=0%2e01&start=",
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=refeicoes-faceis&pmin=0%2e01&start=",
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=congelados&pmin=0%2e01&start=",
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=bebidas&pmin=0%2e01&start=",        // 5k products
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=higiene-beleza&pmin=0%2e01&start=", // 5k products
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=limpeza&pmin=0%2e01&start=",
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=livraria-papelaria&pmin=0%2e01&start=", // 7k products
		// "https://www.continente.pt/on/demandware.store/Sites-continente-Site/default/Search-UpdateGrid?cgid=mercearias&pmin=0%2e01&start=",
	}
)

func main() {
	col := colly.NewCollector(colly.Async(true))
	col.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	col.Limit(&colly.LimitRule{
		Parallelism: 5,
		RandomDelay: 1 * time.Second,
	})

	var mu sync.Mutex
	var wg sync.WaitGroup
	products := make([]ContinenteResponse, 0)

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
		mu.Lock()
		products = append(products, product)
		mu.Unlock()
	})

	col.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	})

	col.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	visitPage := func(baseURL string, pageNumber int, wg *sync.WaitGroup) error {
		defer wg.Done()
		url := fmt.Sprintf("%s%d&sz=24", baseURL, pageNumber*24)
		return col.Visit(url)
	}

	for _, baseURL := range urlsToScrape {
		wg.Add(1)
		go func(baseURL string) {
			defer wg.Done() // Ensure that the wait group counter is decremented
			pageNumber := 0
			for {
				pageWG := &sync.WaitGroup{}
				mu.Lock()
				productCount := len(products)
				mu.Unlock()
				pageWG.Add(1)
				err := visitPage(baseURL, pageNumber, pageWG)
				if err != nil {
					fmt.Println("Error visiting page:", err)
					pageWG.Done() // Ensure we mark the wait group as done if there's an error
					break
				}
				pageWG.Wait() // Wait for the current page request to complete

				if productCount == 0 {
					fmt.Println("No products found on page, ending scraping for category:", baseURL)
					break
				}

				mu.Lock()
				newProductCount := len(products)
				mu.Unlock()

				if newProductCount == productCount {
					fmt.Println("No new products found, ending scraping for category:", baseURL)
					break
				}

				pageNumber++
			}
		}(baseURL)
	}

	wg.Wait()
	col.Wait() // Ensure all async requests are finished

	fmt.Println("TOTAL PRODUCTS LENGTH:", len(products))
}
