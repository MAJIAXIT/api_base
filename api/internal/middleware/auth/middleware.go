package auth

import (
	"net/http"
	"strings"

	"github.com/MAJIAXIT/api_base/api/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	WSAuthMiddleware() gin.HandlerFunc
}

type middleware struct {
	authService auth.Service
}

func New(authService auth.Service) Middleware {
	return &middleware{
		authService: authService,
	}
}

func (m *middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := m.authService.ValidateToken(tokenParts[1], auth.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("authMethod", "jwt")
		c.Set("userID", claims.UserID)
		c.Set("login", claims.Login)

		c.Next()
	}
}

func (m *middleware) WSAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		claims, err := m.authService.ValidateToken(token, auth.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("authMethod", "jwt")
		c.Set("userID", claims.UserID)
		c.Set("login", claims.Login)

		c.Next()
	}
}
