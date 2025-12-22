package repository

import "gorm.io/gorm"

type Repository struct {
	RUser    RUser
	RProduct RProduct
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		RUser:    NewRUser(db),
		RProduct: NewRProduct(db),
	}
}
