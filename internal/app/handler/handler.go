package handler

import (
	"app/internal/app/service"
	"app/pkg/config"
)

type Handler struct {
	Product HProduct
	User    HUser
}

func NewHandler(service *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		Product: NewHProduct(service.SProduct),
		User:    NewHUser(service.SUser, cfg),
	}
}
