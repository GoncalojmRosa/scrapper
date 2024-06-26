package store

import (
	"database/sql"

	"github.com/GoncalojmRosa/scrapper/types"
)

type Storage struct {
	db *sql.DB
}

type Store interface {
	InsertProduct(name, price, img string) error
	GetProducts() ([]types.Product, error)
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}
