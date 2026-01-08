package middleware

import (
	resp "app/internal/dto/response"
	"app/pkg/config"
	"app/pkg/database/redis"
	"app/pkg/toolkit"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/* --------------------------------- Function -------------------------------- */
func Meta() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Set("request_id", uuid.New().String())

		c.Set("response_time", startTime)

		c.Next()
	}
}

func Auth(rds redis.Redis, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			toolkit.ResponseError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// bearer token
		token = strings.TrimPrefix(token, "Bearer ")

		// verify token
		claims, err := VerifyToken(token, cfg.PublicKey)
		if err != nil {
			toolkit.ResponseError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// check redis session
		redisClaims, err := redisSession(c.Request.Context(), rds, claims.JTI)
		if err != nil {
			toolkit.ResponseError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// set context
		c.Set("user_id", redisClaims.ID)
		c.Set("username", redisClaims.Username)
		c.Set("email", redisClaims.Email)

		c.Next()
	}
}

/* --------------------------------- Helper Function -------------------------------- */
func redisSession(ctx context.Context, rds redis.Redis, jti string) (*resp.SessionRecord, error) {
	// check jti in redis
	key := fmt.Sprintf("session:%s", jti)

	// get redis
	value, err := rds.Get(ctx, key)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	// Marshal redis
	var redisClaims resp.SessionRecord
	if err := json.Unmarshal([]byte(value), &redisClaims); err != nil {
		return nil, err
	}

	// extend redis session
	if err := rds.Expire(ctx, key, time.Minute*15); err != nil {
		return nil, err
	}

	return &redisClaims, nil
}
