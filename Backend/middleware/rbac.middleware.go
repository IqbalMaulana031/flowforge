package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"flowforge-api/common/constant"
	"flowforge-api/response"
)

const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		role := c.GetString(RoleKey)
		if _, ok := allowed[role]; !ok {
			response.Error(c, http.StatusForbidden, constant.ErrorCodeForbidden, "role is not allowed")
			c.Abort()
			return
		}
		c.Next()
	}
}
