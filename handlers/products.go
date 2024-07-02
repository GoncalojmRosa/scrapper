package handlers

import (
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
