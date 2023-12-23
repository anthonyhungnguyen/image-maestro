package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get bearer token from header
		authHeader := c.GetHeader("Authorization")

		// Check if bearer token is valid
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token := authHeader[7:]

			// Validate token
			if validateToken(token) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		c.Abort()
	}
}

// TODO: Implement validateToken
func validateToken(token string) bool {
	return true
}
