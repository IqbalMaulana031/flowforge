package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"flowforge-api/config"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	allowAll := false
	for _, origin := range cfg.CORS.AllowedOrigins {
		if origin == "*" {
			allowAll = true
			continue
		}
		allowed[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok || allowAll || len(allowed) == 0 {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", strings.Join([]string{"Authorization", "Content-Type", "X-Request-ID", "X-Tenant-ID"}, ", "))
		c.Header("Access-Control-Allow-Methods", strings.Join([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, ", "))
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
