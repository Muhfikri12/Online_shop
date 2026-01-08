package handler

import (
	"app/internal/app/service"
	"app/internal/dto/request"
	"app/pkg/toolkit"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HUser interface {
	Login(c *gin.Context)
}

type hUser struct {
	sUser service.SUser
}

func NewHUser(sUser service.SUser) HUser {
	return &hUser{
		sUser: sUser,
	}
}

func (h *hUser) Login(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		toolkit.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.sUser.Login(c.Request.Context(), req)
	if err != nil {
		toolkit.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	toolkit.ResponseOK(c, user)
}
