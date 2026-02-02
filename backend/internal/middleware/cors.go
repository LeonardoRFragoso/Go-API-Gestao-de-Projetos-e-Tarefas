package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(origins string) gin.HandlerFunc {
	allowedOrigins := strings.Split(origins, ",")

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		for _, allowed := range allowedOrigins {
			if strings.TrimSpace(allowed) == origin || allowed == "*" {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
