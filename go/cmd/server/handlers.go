package server

import (
	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	validKeys := []string{"SECRET123", "BACKUP456", "TEST789"}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		// Handle missing API key
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "API key required in X-API-Key header",
			})
			return
		}

		// Check if key is valid
		valid := false
		for _, key := range validKeys {
			if apiKey == key {
				valid = true
				break
			}
		}

		if !valid {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid API key",
				"hint":  "Use one of: SECRET123, BACKUP456, TEST789",
			})
			return
		}

		c.Next()
	}
}
