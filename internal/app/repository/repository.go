package repository

import "gorm.io/gorm"

type Repository struct {
	rProduct RProduct
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		rProduct: NewRProduct(db),
	}
}
