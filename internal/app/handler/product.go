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

// FindByUUID godoc
// @Summary      Get product by UUID
// @Description  Returns a single product by its UUID. Requires authentication.
// @Tags         products
// @Produce      json
// @Param        uuid   path      string  true  "Product UUID"
// @Success      200    {object}  toolkit.Response
// @Failure      401    {object}  toolkit.Response
// @Failure      500    {object}  toolkit.Response
// @Security     BearerAuth
// @Router       /product/{uuid} [get]
func (h *hProduct) FindByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	product, err := h.sProduct.FindByUUID(c.Request.Context(), uuid)
	if err != nil {
		toolkit.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	toolkit.ResponseOK(c, product)
}

// FindAll godoc
// @Summary      List products
// @Description  Returns a paginated list of products. Requires authentication.
// @Tags         products
// @Produce      json
// @Success      200  {object}  toolkit.Response
// @Failure      401  {object}  toolkit.Response
// @Failure      500  {object}  toolkit.Response
// @Security     BearerAuth
// @Router       /products [get]
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
