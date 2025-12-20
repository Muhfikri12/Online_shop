package service

import "app/internal/app/repository"

type Service struct {
	SProduct SProduct
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		SProduct: NewSProduct(repo.RProduct),
	}
}
