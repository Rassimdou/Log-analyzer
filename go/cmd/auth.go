package main 


import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	 validKeys := []string{"SECRET887" , "BACKUP444" , "TEST1234"} // Replace with your actual keys
	return func (c *gin.Context){
		apiKey := c.GetHeader("X-API-Key")

		//handle missing key
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key is required",
			})
			return
		}
		
		
		// Check if the API key is valid 
		valid := false
		for _, key := range validKeys {
			if apiKey == key {
				valid = true
				break
			}
		}

		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			return
		}
		c.Next() // Continue to the next handler if the key is valid
	}
}