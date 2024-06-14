package main

import (
	"github.com/gocolly/colly"
)

func main() {
	col := colly.NewCollector()

	col.Visit("https://www.scrapingcourse.com/ecommerce/")

	col.OnHTML("li.product", func(e *colly.HTMLElement) {
		e.ForEach("h2", func(_ int, h *colly.HTMLElement) {
			println(h.Text)
		})
	})

}
