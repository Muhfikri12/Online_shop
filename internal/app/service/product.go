package service

import (
	"app/internal/app/repository"
	resp "app/internal/dto/response"
	"context"
)

type SProduct interface {
	FindByUUID(ctx context.Context, uuid string) (*resp.RespProduct, error)
}

type sProduct struct {
	rProduct repository.RProduct
}

func NewSProduct(rProduct repository.RProduct) SProduct {
	return &sProduct{
		rProduct: rProduct,
	}
}

func (s *sProduct) FindByUUID(ctx context.Context, uuid string) (*resp.RespProduct, error) {
	return s.rProduct.FindByUUID(ctx, uuid)
}
