package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status":  "operational",
			"message": "Security Analyzer Active! ",
		})
	})

	r.POST("/ingest", func(c *gin.Context) {
		//define a struct to bind the incomming JSON data
		type LogEntry struct {
			Source  string `json:"source"`
			Message string `json:"message"`
		}
		var LOG LogEntry
		//bind the JSON to our struct
		if err := c.ShouldBindJSON(&LOG); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request format"})
			return
		}

		println("Received log from", LOG.Source+":", LOG.Message)

		c.JSON(200, gin.H{"status": "received"})

	})


	// Create rate limiter: 10 requests/minute per IP
	limiter := NewRateLimiter(10, time.Minute)

	// Apply globally or to specific routes
	r.Use(limiter.Middleware())
	
	r.POST("/ingest", APIKeyAuthMiddleware(), IngestHandler)

	r.Run()

}
