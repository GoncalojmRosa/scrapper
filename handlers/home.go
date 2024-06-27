package handlers

import (
	"net/http"

	"github.com/GoncalojmRosa/scrapper/types"
	"github.com/GoncalojmRosa/scrapper/views"
)

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	var ProductList []types.Product = []types.Product{
		{
			Name:  "Product 1",
			Price: "10.00",
			Img:   "https://via.placeholder.com/150",
		},
		{
			Name:  "Product 2",
			Price: "20.00",
			Img:   "https://via.placeholder.com/150",
		},
	}

	views.Home(ProductList).Render(r.Context(), w)
}
