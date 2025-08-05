package mw

import (
	"fmt"
	"net/http"
	"sync"
	"xi/app/lib"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMw struct {
	jwtSecret []byte

	once sync.Once
	mu   sync.RWMutex
}

// Global singleton instance
var Auth = &AuthMw{
	jwtSecret: []byte(lib.Env.Get("JWT_SECRET")),
}

// Middleware to verify token
func (a *AuthMw) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		// Expected format: "Bearer <token>"
		var tokenStr string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
		if err != nil || tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return a.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}
