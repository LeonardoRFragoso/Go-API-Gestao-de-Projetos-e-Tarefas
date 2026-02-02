package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taskmanager/backend/internal/handler"
	"github.com/taskmanager/backend/internal/service"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			handler.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
