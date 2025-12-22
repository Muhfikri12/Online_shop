package handler

import (
	"app/internal/app/service"
)

type Handler struct {
	Product HProduct
	User    HUser
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Product: NewHProduct(service.SProduct),
		User:    NewHUser(service.SUser),
	}
}
