package handler

import (
	"app/internal/app/service"
	"app/pkg/toolkit"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HProduct interface {
	FindByUUID(c *gin.Context)
	FindAll(c *gin.Context)
}

type hProduct struct {
	sProduct service.SProduct
}

func NewHProduct(sProduct service.SProduct) HProduct {
	return &hProduct{
		sProduct: sProduct,
	}
}

func (h *hProduct) FindByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	product, err := h.sProduct.FindByUUID(c.Request.Context(), uuid)
	if err != nil {
		toolkit.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	toolkit.ResponseOK(c, product)
}

func (h *hProduct) FindAll(c *gin.Context) {
	products, err := h.sProduct.FindAll(c.Request.Context())
	if err != nil {
		toolkit.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("username", c.GetString("username"))

	toolkit.ResponsePage(c, products, toolkit.Pagination{
		Page:  1,
		Limit: 10,
	})
}
