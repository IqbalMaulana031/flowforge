package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"flowforge-api/common/constant"
	"flowforge-api/response"
)

const TenantIDKey = "tenantId"

func TenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(TenantIDKey) == "" {
			response.Error(c, http.StatusUnauthorized, constant.ErrorCodeUnauthorized, "tenant context is missing")
			c.Abort()
			return
		}
		c.Next()
	}
}
