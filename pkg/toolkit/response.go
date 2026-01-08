package toolkit

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Meta    Meta        `json:"meta"`
	Data    interface{} `json:"data"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

type Meta struct {
	Pagination     *Pagination `json:"pagination,omitempty"`
	RequestID      string      `json:"request_id"`
	ResponseTimeMS string      `json:"response_time_ms"`
	Timestamp      time.Time   `json:"timestamp"`
}

func ResponseOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success",
		Meta:    buildMeta(c),
		Data:    data,
	})
}

func ResponseError(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Status:  "error",
		Message: message,
		Meta:    buildMeta(c),
		Data:    nil,
	})
}

func ResponsePage(c *gin.Context, data interface{}, page Pagination) {
	meta := buildMeta(c)
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success",
		Meta: Meta{
			Pagination:     &page,
			RequestID:      meta.RequestID,
			ResponseTimeMS: meta.ResponseTimeMS,
			Timestamp:      meta.Timestamp,
		},
		Data: data,
	})
}

func buildMeta(c *gin.Context) Meta {
	var rt int64
	if d, ok := c.Get("response_time"); ok {
		if startTime, ok := d.(time.Time); ok {
			rt = time.Since(startTime).Milliseconds()
		}
	}

	return Meta{
		RequestID:      c.GetString("request_id"),
		ResponseTimeMS: fmt.Sprintf("%d ms", rt),
		Timestamp:      time.Now(),
	}
}
