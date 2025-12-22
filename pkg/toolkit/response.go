package toolkit

import (
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
	ResponseTimeMS int64       `json:"response_time_ms"`
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
	var responseTime int64
	if startTime, exists := c.Get("start_time"); exists {
		if t, ok := startTime.(time.Time); ok {
			responseTime = time.Since(t).Milliseconds()
		}
	}

	return Meta{
		RequestID:      c.GetString("request_id"),
		ResponseTimeMS: responseTime,
		Timestamp:      time.Now(),
	}
}
