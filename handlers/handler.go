package handlers

import "github.com/GoncalojmRosa/scrapper/store"

type Handler struct {
	store *store.Storage
}

func New(store *store.Storage) *Handler {
	return &Handler{
		store: store,
	}
}
