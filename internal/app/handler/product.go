package handler

import (
	"app/internal/app/service"
	resp "app/internal/dto/response"
	"context"
)

type HProduct interface {
	FindByUUID(ctx context.Context, uuid string) (*resp.RespProduct, error)
}

type hProduct struct {
	sProduct service.SProduct
}

func NewHProduct(sProduct service.SProduct) HProduct {
	return &hProduct{
		sProduct: sProduct,
	}
}

func (h *hProduct) FindByUUID(ctx context.Context, uuid string) (*resp.RespProduct, error) {
	return h.sProduct.FindByUUID(ctx, uuid)
}
