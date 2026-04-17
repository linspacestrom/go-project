package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const allowAnyOrigin = "*"

func CORS(allowedOrigins []string) gin.HandlerFunc {
	origins := normalizeOrigins(allowedOrigins)
	allowAny := origins[allowAnyOrigin]

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		c.Header("Vary", "Origin")
		c.Header("Vary", "Access-Control-Request-Method")
		c.Header("Vary", "Access-Control-Request-Headers")

		if allowAny || origins[origin] {
			if allowAny {
				c.Header("Access-Control-Allow-Origin", allowAnyOrigin)
			} else {
				c.Header("Access-Control-Allow-Origin", origin)
			}
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Request-ID")
			c.Header("Access-Control-Max-Age", "600")
		}

		if c.Request.Method == http.MethodOptions {
			if !(allowAny || origins[origin]) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func normalizeOrigins(origins []string) map[string]bool {
	result := make(map[string]bool, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "" {
			continue
		}
		result[trimmed] = true
	}

	return result
}
