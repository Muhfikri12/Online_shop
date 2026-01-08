package middleware

import (
	"app/pkg/toolkit"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Meta() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Set("request_id", uuid.New().String())

		c.Set("response_time", startTime)

		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			toolkit.ResponseError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

	}
}
