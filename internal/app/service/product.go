package service

import (
	"app/internal/app/repository"
	resp "app/internal/dto/response"
	"app/internal/model"
	"context"
)

/* --------------------------------- Interface --------------------------------- */
type SProduct interface {
	FindByUUID(ctx context.Context, uuid string) (*resp.RespProduct, error)
	FindAll(ctx context.Context) ([]model.Product, error)
}

type sProduct struct {
	rProduct repository.RProduct
}

func NewSProduct(rProduct repository.RProduct) SProduct {
	return &sProduct{
		rProduct: rProduct,
	}
}

/* --------------------------------- Function -------------------------------- */
func (s *sProduct) FindByUUID(ctx context.Context, uuid string) (*resp.RespProduct, error) {
	return s.rProduct.FindByUUID(ctx, uuid)
}

func (s *sProduct) FindAll(ctx context.Context) ([]model.Product, error) {
	return s.rProduct.FindAll(ctx)
}
