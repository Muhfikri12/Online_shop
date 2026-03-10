package service

import (
	"app/internal/app/repository"
	"app/pkg/config"
	rds "app/pkg/database/redis"
)

type Service struct {
	SProduct SProduct
	SUser    SUser
}

func NewService(repo *repository.Repository, cfg *config.Config, rds rds.Redis) *Service {
	return &Service{
		SProduct: NewSProduct(repo.RProduct),
		SUser:    NewSUser(repo.RUser, repo.RAuth, cfg, rds),
	}
}
