package handlers

import (
	"fmt"
	"net/http"

	"github.com/GoncalojmRosa/scrapper/views"
)

func (h *Handler) HandleListProducts(w http.ResponseWriter, r *http.Request) {

	prods, err := h.store.GetProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	views.Products(prods).Render(r.Context(), w)
}

func (h *Handler) HandleProductSearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Searching for product:", r.FormValue("name"))
	prods, err := h.store.GetProductByName(r.FormValue("name"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	views.Products(prods).Render(r.Context(), w)
}
