package repository

import "gorm.io/gorm"

type Repository struct {
	RProduct RProduct
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		RProduct: NewRProduct(db),
	}
}
