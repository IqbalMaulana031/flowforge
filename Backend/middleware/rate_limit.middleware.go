package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"flowforge-api/config"
	"flowforge-api/response"
)

func RateLimit(cfg *config.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}
		limit := cfg.RateLimit.Anonymous
		key := "rate:anon:" + c.ClientIP()
		if tenantID := c.GetString(TenantIDKey); tenantID != "" {
			limit = cfg.RateLimit.Authenticated
			key = "rate:tenant:" + tenantID
		}
		ctx := context.Background()
		count, err := redisClient.Incr(ctx, key).Result()
		if err == nil && count == 1 {
			_ = redisClient.Expire(ctx, key, time.Minute).Err()
		}
		if err == nil && count > int64(limit) {
			response.Error(c, http.StatusTooManyRequests, "RATE_LIMITED", "rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}
