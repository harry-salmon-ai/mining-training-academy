package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		// Check Authorization header first
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Fall back to cookie
		if tokenString == "" {
			tokenString, _ = c.Cookie("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		c.Set("siteLocation", claims.SiteLocation)
		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
		if tokenString == "" {
			tokenString, _ = c.Cookie("token")
		}

		if tokenString != "" {
			claims, err := ValidateToken(tokenString)
			if err == nil {
				c.Set("userId", claims.UserID)
				c.Set("userEmail", claims.Email)
				c.Set("userRole", claims.Role)
				c.Set("siteLocation", claims.SiteLocation)
			}
		}

		c.Next()
	}
}
