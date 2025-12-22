package service

import (
	"app/internal/app/repository"
	"app/pkg/config"
)

type Service struct {
	SProduct SProduct
	SUser    SUser
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		SProduct: NewSProduct(repo.RProduct),
		SUser:    NewSUser(repo.RUser, cfg),
	}
}
