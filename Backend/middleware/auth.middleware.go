package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"flowforge-api/common/constant"
	"flowforge-api/config"
	"flowforge-api/response"
	"flowforge-api/utils"
)

const (
	UserIDKey      = "userId"
	RoleKey        = "role"
	AccessTokenKey = "accessToken"
)

func Auth(cfg *config.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, constant.ErrorCodeUnauthorized, "missing bearer token")
			c.Abort()
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		claims, err := utils.ParseToken(cfg, tokenString)
		if err != nil || claims.TokenUse != "access" {
			response.Error(c, http.StatusUnauthorized, constant.ErrorCodeUnauthorized, "invalid bearer token")
			c.Abort()
			return
		}
		if redisClient != nil {
			if exists, _ := redisClient.Exists(context.Background(), "jwt:blacklist:"+tokenString).Result(); exists > 0 {
				response.Error(c, http.StatusUnauthorized, constant.ErrorCodeUnauthorized, "token has been revoked")
				c.Abort()
				return
			}
		}
		c.Set(AccessTokenKey, tokenString)
		c.Set(UserIDKey, claims.UserID)
		c.Set(TenantIDKey, claims.TenantID)
		c.Set(RoleKey, claims.Role)
		c.Next()
	}
}
